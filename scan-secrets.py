#!/usr/bin/env python3
"""
Secret Scanner — Detect leaked credentials in your codebase.
Usage:
    python scan-secrets.py                          # Scan all staged files (pre-commit mode)
    python scan-secrets.py --path /path/to/repo     # Scan entire repo
    python scan-secrets.py --path /path/to/file.py  # Scan single file
    python scan-secrets.py --install-hook           # Install as git pre-commit hook
"""

import os
import re
import sys
import json
import subprocess
from pathlib import Path

# ─── COLOR ───
RED = "\033[91m"
GREEN = "\033[92m"
YELLOW = "\033[93m"
CYAN = "\033[96m"
BOLD = "\033[1m"
RESET = "\033[0m"

# ─── RULES ───
# Semua pattern untuk detect secret di codebase
SECRET_RULES = [
    # OpenAI / LLM API Keys
    (r'sk-[a-zA-Z0-9]{20,}', "OpenAI API Key (sk-...)"),
    (r'sk-ant-[a-zA-Z0-9]{20,}', "Anthropic API Key (sk-ant-...)"),
    
    # AWS
    (r'AKIA[0-9A-Z]{16}', "AWS Access Key ID (AKIA...)"),
    (r'(?i)aws.*secret.*access.*key["\']?\s*[:=]\s*["\']?[A-Za-z0-9/+=]{40}["\']?', "AWS Secret Access Key"),
    
    # Google Cloud
    (r'AIza[0-9A-Za-z\-_]{35}', "Google API Key (AIza...)"),
    (r'(?i)"type"\s*:\s*"service_account"', "Google Service Account JSON"),
    
    # GitHub
    (r'ghp_[a-zA-Z0-9]{36}', "GitHub Personal Access Token (ghp_...)"),
    (r'gho_[a-zA-Z0-9]{36}', "GitHub OAuth Access Token (gho_...)"),
    (r'github_pat_[a-zA-Z0-9]{22,}', "GitHub Fine-Grained PAT"),
    
    # Slack
    (r'xoxb-[0-9]{10,13}-[0-9]{10,13}-[a-zA-Z0-9]{24}', "Slack Bot Token (xoxb-...)"),
    (r'xapp-[0-9]{10,13}-[a-zA-Z0-9]{24}', "Slack App Token (xapp-...)"),
    
    # Telegram
    (r'[0-9]{8,10}:[A-Za-z0-9_-]{35}', "Telegram Bot Token"),
    
    # JWT / Tokens
    (r'eyJ[a-zA-Z0-9_-]{10,}\.[a-zA-Z0-9_-]{10,}\.[a-zA-Z0-9_-]{10,}', "JWT Token"),
    
    # Private Keys
    (r'-----BEGIN\s?RSA\s?PRIVATE\s?KEY-----', "RSA Private Key"),
    (r'-----BEGIN\s?EC\s?PRIVATE\s?KEY-----', "EC Private Key"),
    (r'-----BEGIN\s?OPENSSH\s?PRIVATE\s?KEY-----', "OpenSSH Private Key"),
    
    # Database URLs
    (r'postgresql://[^:@]+:[^@]+@', "PostgreSQL Connection String (with password)"),
    (r'mongodb(?:\+srv)?://[^:]+:[^@]+@', "MongoDB Connection String (with password)"),
    (r'mysql://[^:@]+:[^@]+@', "MySQL Connection String (with password)"),
    (r'redis://[^:@]+:[^@]+@', "Redis Connection String (with password)"),
    
    # Generic Password in connection
    (r'(?i)(password|passwd|pwd)\s*[=:]\s*["\'][^"\'\s]+["\']', "Generic Password Assignment"),
    
    # Heroku
    (r'heroku[a-z]{0,10}:\/\/[a-z0-9]{20,}:', "Heroku API Key"),
    
    # Discord
    (r'(?:discord|discordapp)\.com\/api\/webhooks\/[0-9]+\/[a-zA-Z0-9_-]+', "Discord Webhook URL"),
    
    # Stripe
    (r'(?:sk|pk)_(?:live|test)_[a-zA-Z0-9]{24,}', "Stripe API Key"),
    
    # npm / package registry
    (r'(?i)npm_token["\']?\s*[:=]\s*["\']?[a-zA-Z0-9]{20,}["\']?', "npm Token"),
    
    # .env file contents pattern
    (r'(?i)(api_key|api_secret|secret_key|access_key|private_key)\s*=\s*["\']?[a-zA-Z0-9_\-/.+=]{8,}["\']?', "Potential Secret in .env pattern"),
]

# File extensions/names to ALWAYS scan
ALWAYS_SCAN = {'.py', '.js', '.ts', '.jsx', '.tsx', '.go', '.rs', '.rb', '.php', '.java', '.kt', '.sh', '.bash', '.zsh', '.yaml', '.yml', '.toml', '.ini', '.cfg', '.conf', '.env', '.env.example'}

# Files/patterns to NEVER scan (binary, vendor, etc.)
IGNORE_PATTERNS = [
    'node_modules/', '.git/', 'vendor/', 'dist/', 'build/', '.next/',
    '__pycache__/', '*.pyc', '*.pyo', '*.so', '*.dll', '*.dylib', '*.exe',
    '*.png', '*.jpg', '*.jpeg', '*.gif', '*.svg', '*.ico', '*.woff', '*.woff2',
    '*.ttf', '*.eot', '*.zip', '*.tar.gz', '*.tgz', '*.7z', '*.rar',
    '.gitkeep', 'package-lock.json', 'yarn.lock',
]


def should_ignore(filepath: str) -> bool:
    """Cek apakah file perlu di-skip."""
    for pat in IGNORE_PATTERNS:
        if pat.startswith('*.'):
            if filepath.endswith(pat[1:]):
                return True
        elif pat in filepath:
            return True
    return False


def scan_file(filepath: str) -> list[dict]:
    """Scan satu file dan return list of findings."""
    results = []
    
    # Skip .env.example files — ini sengaja placeholder
    if filepath.endswith('.env.example'):
        return results
        
    try:
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
    except Exception:
        return results

    lines = content.split('\n')
    for lineno, line in enumerate(lines, 1):
        for pattern, description in SECRET_RULES:
            matches = re.findall(pattern, line)
            for match in matches:
                # Redact the secret — tampilkan cuma 4 char pertama & terakhir
                secret = str(match)
                if len(secret) > 12:
                    redacted = secret[:4] + "****" + secret[-4:]
                else:
                    redacted = secret
                results.append({
                    'file': filepath,
                    'line': lineno,
                    'description': description,
                    'match': redacted,
                })
    return results


def get_staged_files() -> list[str]:
    """Ambil daftar file yang di-stage untuk di-commit."""
    try:
        result = subprocess.run(
            ['git', 'diff', '--cached', '--name-only', '--diff-filter=ACM'],
            capture_output=True, text=True, check=True
        )
        files = [f.strip() for f in result.stdout.split('\n') if f.strip()]
        return [f for f in files if not should_ignore(f)]
    except (subprocess.CalledProcessError, FileNotFoundError):
        return []


def scan_directory(path: str) -> list[dict]:
    """Scan semua file dalam direktori."""
    results = []
    for root, dirs, files in os.walk(path):
        # Skip direktori yang di-ignore
        dirs[:] = [d for d in dirs if d not in {'.git', 'node_modules', 'vendor', 'dist', 'build', '.next', '__pycache__'}]
        for f in files:
            filepath = os.path.join(root, f)
            relpath = os.path.relpath(filepath, path)
            if not should_ignore(relpath):
                results.extend(scan_file(filepath))
    return results


def print_report(results: list[dict]) -> bool:
    """Print hasil scan. Return True kalo ada yang ketahuan."""
    if not results:
        print(f"\n{GREEN}✅ No secrets found!{' ' * 50}{RESET}")
        return False

    # Group by file
    by_file = {}
    for r in results:
        by_file.setdefault(r['file'], []).append(r)

    print(f"\n{RED}{BOLD}⚠️  WARNING: {len(results)} potential secrets detected!{RESET}")
    print("=" * 70)

    for filepath, findings in sorted(by_file.items()):
        print(f"\n{YELLOW}📁 {filepath}{RESET}")
        for f in findings:
            print(f"   {CYAN}L{f['line']:>4}{RESET}  {RED}🔴{RESET}  {f['description']}")
            print(f"         {YELLOW}Match: {f['match']}{RESET}")

    print()
    print(f"{BOLD}Summary:{RESET}")
    print(f"  {RED}Total secrets found: {len(results)}{RESET}")
    print(f"  Files affected: {len(by_file)}")
    print()
    return True


def install_precommit_hook():
    """Install hook ke .git/hooks/pre-commit."""
    repo_root = subprocess.run(
        ['git', 'rev-parse', '--show-toplevel'],
        capture_output=True, text=True
    ).stdout.strip()

    if not repo_root:
        print(f"{RED}❌ Not inside a Git repository.{RESET}")
        sys.exit(1)

    hooks_dir = os.path.join(repo_root, '.git', 'hooks')
    hook_path = os.path.join(hooks_dir, 'pre-commit')

    hook_script = """#!/bin/bash
# Secret Scanner — Auto-detect leaked credentials before commit
# Diinstall otomatis oleh scan-secrets.py

echo ""
echo "🔍 Scanning staged files for secrets..."
python3 "$(dirname "$0")/../../scan-secrets.py"
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    echo ""
    echo "❌ Commit BLOCKED — potential secrets detected!"
    echo "   Fix them or use 'git commit --no-verify' to bypass (not recommended)."
    exit 1
fi

exit 0
"""

    with open(hook_path, 'w') as f:
        f.write(hook_script)
    os.chmod(hook_path, 0o755)

    # Copy script ke root repo kalo belum ada
    script_src = os.path.abspath(__file__)
    script_dst = os.path.join(repo_root, 'scan-secrets.py')
    if not os.path.exists(script_dst):
        with open(script_src, 'r') as src, open(script_dst, 'w') as dst:
            dst.write(src.read())
        os.chmod(script_dst, 0o755)
        print(f"{GREEN}✅ Copied scan-secrets.py to {script_dst}{RESET}")

    print(f"{GREEN}✅ Pre-commit hook installed at {hook_path}{RESET}")
    print(f"{GREEN}✅ Secrets will be checked automatically before every commit!{RESET}")


def main():
    import argparse
    parser = argparse.ArgumentParser(description='Secret Scanner — Detect leaked credentials')
    parser.add_argument('--path', help='Scan specific path (file or directory)')
    parser.add_argument('--install-hook', action='store_true', help='Install as git pre-commit hook')
    parser.add_argument('--json', action='store_true', help='Output as JSON')
    args = parser.parse_args()

    if args.install_hook:
        install_precommit_hook()
        return

    # Determine files to scan
    if args.path:
        path = args.path
        if os.path.isfile(path):
            results = scan_file(path)
        elif os.path.isdir(path):
            results = scan_directory(path)
        else:
            print(f"{RED}❌ Path not found: {path}{RESET}")
            sys.exit(1)
    else:
        # Default: scan staged files (pre-commit mode)
        staged = get_staged_files()
        if not staged:
            print(f"{YELLOW}ℹ️  No staged files to scan.{RESET}")
            sys.exit(0)
        results = []
        for f in staged:
            results.extend(scan_file(f))

    if args.json:
        print(json.dumps(results, indent=2))
    else:
        has_secrets = print_report(results)
        sys.exit(1 if has_secrets else 0)


if __name__ == '__main__':
    main()

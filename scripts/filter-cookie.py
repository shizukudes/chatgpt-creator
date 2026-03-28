#!/usr/bin/env python3
import json
import sys

def filter_cookies(input_file, output_file):
    with open(input_file, 'r') as f:
        cookies = json.load(f)

    filtered = []
    for c in cookies:
        name = c.get('name', '')
        domain = c.get('domain', '')
        # Only chatgpt.com and .chatgpt.com domains
        if 'chatgpt.com' not in domain:
            continue
        # Only these specific cookies
        if name in ['oai-did', '_cfuvid', '__cf_bm', '__Secure-next-auth.session-token']:
            filtered.append({
                'name': c.get('name'),
                'value': c.get('value'),
                'domain': c.get('domain'),
                'path': c.get('path', '/'),
                'secure': c.get('secure', False),
                'httpOnly': c.get('httpOnly', False)
            })

    with open(output_file, 'w') as f:
        json.dump(filtered, f, indent=2)

    print(f"Filtered {len(filtered)} cookies -> {output_file}")

if __name__ == '__main__':
    input_file = sys.argv[1] if len(sys.argv) > 1 else 'cookies.json'
    output_file = sys.argv[2] if len(sys.argv) > 2 else 'filtered-cookies.json'
    filter_cookies(input_file, output_file)

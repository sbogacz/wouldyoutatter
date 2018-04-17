#!/usr/bin/env python3

""" This is a horrible hack but loading an entire SVG
    (or even XML) parser would be a tragedy. """

import re
import sys

def main():
    file_name = sys.argv[1]
    with open(file_name, 'r') as in_svg:
        svg_body = in_svg.read()
    svg_body = re.sub(r'<text.*text>', '', svg_body)
    with open(file_name, 'w') as out_svg:
        out_svg.write(svg_body)

if __name__ == '__main__':
    main()

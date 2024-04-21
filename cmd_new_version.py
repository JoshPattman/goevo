# This simple script generates a shell script to push a new version of goevo and all of its submodules
# I prefer it when all of the submodules share the same version as the parent module, even if they have no changes

import os
import sys
import subprocess
import re

try:
    version = os.environ['GOEVO_VERSION']
    print(f"Creating new version {version}...")
except KeyError:
    print("GOEVO_VERSION is not set.")
    sys.exit(1)

submodules = []
with open("go.work", "r") as f:
    lines = f.readlines()
    for line in lines:
        line = line.strip()
        if re.match(r'./', line):
            name = line.strip().rstrip("/")
            if name != ".":
                submodules.append(name)

print(f"Submodules: {submodules}")

current_dir = os.getcwd()
print(f"Current directory: {current_dir}")

with open("cmd_new_version.sh", "w") as script:
    script.write(f"#!/bin/bash\n")

    script.write("\n# Ensure all current changes are pushed\n")
    script.write("git add .\n")
    script.write(f"git commit -m 'Bump parent version to {version}' --allow-empty\n")
    script.write("git push\n")

    script.write("\n# Create a new tag and upload it\n")
    script.write(f"git tag -a {version} -m 'Release (goevo) {version}'\n")
    script.write("git push --tags\n")

    script.write("\n# Upgrade all submodules to newest version\n")
    for submodule in submodules:
        script.write(f"cd {submodule}\n")
        script.write(f"go get github.com/JoshPattman/goevo@{version}\n")
        script.write("go mod tidy\n")
        script.write(f"cd {current_dir}\n")

    script.write("\n# Push submodule version changes\n")
    script.write("git add .\n")
    script.write(f"git commit -m 'Bump submodule versions to {version}' --allow-empty\n")
    script.write("git push\n")

    script.write("\n# Tag versions of all submodules\n")
    for submodule in submodules:
        tag_name = submodule.strip("./")
        script.write(f"""git tag -a {tag_name}/{version} -m "Release {version} for {tag_name}"\n""")
    script.write("git push --tags\n")


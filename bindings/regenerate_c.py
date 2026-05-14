#!/usr/bin/env python3
from pathlib import Path
import itertools


def main():
    script_dir = Path(__file__).parent
    repo_dir = script_dir.parent
    dst_dir = script_dir / "c"
    assert dst_dir.exists(), f"{dst_dir} missing"

    # Delete all the old C files
    for f in dst_dir.glob("*.c"):
        f.unlink()

    links = {}
    for src in itertools.chain(
        repo_dir.glob("TPMCmd/tpm/src/**/*.c"),
        repo_dir.glob("TPMCmd/tpm/cryptolibs/Ossl/**/*.c"),
        repo_dir.glob("TPMCmd/tpm/cryptolibs/TpmBigNum/**/*.c"),
        repo_dir.glob("TPMCmd/TpmConfiguration/**/*.c"),
    ):
        # Trim any leading underscores to stop CGO from ignoring the file.
        dst = dst_dir / src.name.lstrip("_")
        assert dst not in links, f"Name collision: {links[dst]} and {src}"

        links[dst] = src
        with open(dst, "w") as f:
            f.write(f'#include "../../{src.relative_to(repo_dir)}"\n')

    print(f"Created {len(links)} wrapper C files in {dst_dir}")


if __name__ == "__main__":
    main()

# `bump`: Blind Updater for Deno

Update yer modules really darn fast.

## Usage

With directly-linked modules in a file:

```
$ cat test/sample.ts
export { BumperService } from "https://raw.githubusercontent.com/drashland/services/v0.2.5/ci/bumper_service.ts";
export * as Line from "https://deno.land/x/line@v1.0.1/mod.ts";
import { ConsoleLogger } from "https://deno.land/x/unilogger@v1.0.4/mod.ts";
const consoleLogger = new ConsoleLogger({});
export { consoleLogger as ConsoleLogger };
export { assertEquals } from "https://deno.land/std@0.157.0/testing/asserts.ts";
export * as colours from "https://deno.land/std@0.157.0/fmt/colors.ts";
$ bump test/
Analyzing test/ ...
Found source file: sample.ts
Updated 5 modules in test/sample.ts
$ cat test/sample.ts
export { BumperService } from "https://raw.githubusercontent.com/drashland/services/v0.2.5/ci/bumper_service.ts";
export * as Line from "https://deno.land/x/line@v1.0.1/mod.ts";
import { ConsoleLogger } from "https://deno.land/x/unilogger@v1.1.0/mod.ts";
const consoleLogger = new ConsoleLogger({});
export { consoleLogger as ConsoleLogger };
export { assertEquals } from "https://deno.land/std@0.160.0/testing/asserts.ts";
export * as colours from "https://deno.land/std@0.160.0/fmt/colors.ts";
```

With an `import_map`:

```
$ cat test/import_map.json
{
  "imports": {
    "$fresh/": "../",

    "twind": "https://esm.sh/twind@0.16.17",
    "twind/": "https://esm.sh/twind@0.16.17/",

    "preact": "https://esm.sh/preact@10.11.0",
    "preact/": "https://esm.sh/preact@10.11.0/",
    "preact-render-to-string": "https://esm.sh/*preact-render-to-string@5.2.4",
    "@preact/signals": "https://esm.sh/*@preact/signals@1.0.3",
    "@preact/signals-core": "https://esm.sh/@preact/signals-core@1.0.1",

    "$std/": "https://deno.land/std@0.150.0/"
  }
}
$ bump test/
Analyzing test/ ...
Found source file: import_map.json
Updated 8 modules in test/import_map.json
$ cat test/import_map.json
{
  "imports": {
    "$fresh/": "../",

    "twind": "https://esm.sh/twind@0.16.17",
    "twind/": "https://esm.sh/twind@0.16.17/",

    "preact": "https://esm.sh/preact@10.11.2",
    "preact/": "https://esm.sh/preact@10.11.2/",
    "preact-render-to-string": "https://esm.sh/*preact-render-to-string@5.2.6",
    "@preact/signals": "https://esm.sh/*@preact/signals@1.1.2",
    "@preact/signals-core": "https://esm.sh/@preact/signals-core@1.2.2",

    "$std/": "https://deno.land/std@0.160.0/"
  }
}
```

## Installation
```
go install github.com/bit-bandit/bump
```

## Notes

- Current only works well with `deno.land` and `esm.sh`. Most other repositories
  are untested, and may not work as intended.
- This will only update specifically versioned modules. Unversioned modules will
  be ignored.
- This is experimental software. Use at your own risk.

## License

0BSD

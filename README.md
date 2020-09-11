
# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

# Building the Binaries

Binaries will be placed under `bin/` directory. The name would be `GuestAgentTestExtension_windows.exe` and `GuestAgentTestExtension_linux` based on the OS you build it for.
To build the binaries, run either of the following commands -

- Build Without dependencies -

  - `make build_all` - Use this command to build windows and linux binaries.
  - `make build_windows` - Build only windows binaries.
  - `make build_linux` - Build only linux binaries.

- Build with dependencies -
  - `make build_all_with_deps` - Use this command to build windows and linux binaries with dependencies.
  - `make build_windows` - Build only windows binaries with dependencies.
  - `make build_linux` - Build only linux binaries with dependencies.

- Download dependencies -
  - `make deps`

- Clean project - Delete all binaries
  - `make clean`

# Running the program

To run the program, just run the binary.

For usage information and a list of all available options, run the binary with the help flag ie (--h or --help).

```text
.\bin\GuestAgentTestExtension_win.exe --help

Usage of C:\GIT-GA-Test-Extension\guest-agent-test-extension\bin\GuestAgentTestExtension_win.exe:
  -command string
        Valid commands are install, enable, update, disable and uninstall. Usage: --command=install
  -failCommandFile string
        Path to the JSON file loction. Usage --failCommandFile="test.json"
```

Example:
`.\bin\GuestAgentTestExtension_win.exe -command=install --failCommandFile=="sample-fail-commands.json"`

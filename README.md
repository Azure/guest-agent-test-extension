
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

To run the program, you should use theprovided makefile commands. This will create a "bin" directory and then put
executables for all supported operating systems into the folder.

For usage information, run the binary with the help flag ie (--h or --help)

# Sample command to load onto a live VM

Set-AzVMExtension -ResourceGroupName $rgName -Location "centraluseuap" -VMName $vmName -Name "GATestExtension" -Publisher "Microsoft.Azure.Extensions.Edp" -Type "GATestExtGo" -TypeHandlerVersion "1.0" -Settings $settings

# Runtime Configuration 

This file controls the behavior of the extension once it is compiled. Currently this just takes the form of failCommand specifciations

FailCommands is essentially a list of commands that we want to fail. We can control the command to fail, the error Message at failure, the exit code, and whether we want to report the error status correctly. If reportStatusCorrectly is false, the extension status will be transitioning at the end of the execution.

Sample:
{
    "failCommands": [
        {
            "command" : "install",
            "errorMessage" : "install failure due to specification in failCommand",
            "exitCode" : "9",
            "reportStatusCorrectly" : "false"
        },
        {
            "command" : "enable",
            "errorMessage" : "enable failure due to specification in failCommand",
            "exitCode" : "10",
            "reportStatusCorrectly" : "true"
        }
    ]
}

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

To run the program, you should use the make build_all command. This will create a "bin" directory and then put
executables for all supported operating systems into the folder.

For usage information, run the binary with the help flag ie (--h or --help)

# Sample command to load onto a live VM

Set-AzVMExtension -ResourceGroupName $rgName -Location "centraluseuap" -VMName $vmName -Name "GATestExtension" -Publisher "Microsoft.Azure.Extensions.Edp" -Type "GATestExtGo" -TypeHandlerVersion "1.0" -Settings $settings
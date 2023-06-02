# NymSTR-CLI

NymSTR-CLI is a comprehensive privacy tool built on top of the Nym mixnet. With a focus on privacy and anonymity, it offers a suite of applications to cater to various needs while ensuring that your data remains secure and private. Its functionalities extend to email, file sharing, and wallet operations.

## Features

* **Email**: Send and receive encrypted emails over the Nym mixnet.
* **File Sharing**: Share files with confidence knowing they are encrypted and transferred over the secure mixnet.
* **Wallet**: Conduct transactions securely with the wallet feature.

## Requirements

* Go 1.x
* Nym Client

## How to Run

1. Initialize the Nym Client with your desired ID. You can do this by running `nym-client init --id nymstr-cli`. If you want to use a different ID, make sure to change it in the code as well.

    ```bash
    nym-client init --id nymstr-cli
    ```

2. Navigate to the directory containing the project's Go files and run the program:

    ```bash
    go run main.go
    ```

## Upcoming Features (TODO)

* Decrypt incoming messages
* Remove pop-up contact list
* Write contacts to a file in .nym/ directory
* Properly run Nym Client in the background
* Keep `main.go` as the main entry point, and reorganize other packages into their proper place

## Contributing

We welcome contributions to this project. Please follow the [contribution guidelines](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).


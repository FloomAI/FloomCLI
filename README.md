
# Floom CLI



Floom CLI is a powerful command-line interface designed to simplify the configuration, management, and deployment of Floom environments. It allows developers and DevOps professionals to efficiently manage Docker-compose based Floom environments and deploy YAML configuration files for customizing Floom instances.



## Features



-  **Environment Management**: Start, stop, and manage your Floom environments with ease.

-  **Configuration Deployment**: Easily deploy configuration files to tailor Floom instances to your needs.

-  **Scalability**: Seamlessly scale your AI-powered applications.

-  **Integration**: Designed to fit into existing development workflows, offering a robust CLI toolset for Floom environment management.



## Installation



### Prerequisites



- Docker and Docker Compose

- Go version 1.15 or higher (for building from source)



### From Binaries



Download the latest release for your platform from the [Releases](#) page and extract the binary to a location in your system's PATH.



### Building From Source



To build Floom CLI from source, clone the repository and use Go to compile:



```bash

git  clone  https://github.com/FloomAI/FloomCLI.git

cd  FloomCLI

go  build  -o  floom  .

```



## Usage



Here's how you can use Floom CLI to manage your environments:



### Start an Environment



```bash

floom  start

```



### Stop an Environment



```bash

floom  stop

```



### Deploy a Configuration



```bash

floom  deploy [cloud/local/endpoint] path/to/config.yml

```



For more detailed information on commands and their usage, run:



```bash

floom  --help

```



## Contributing



We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for more information on how to get started.



## License



Floom CLI is open-sourced under the MIT License. See the [LICENSE](LICENSE) file for more details.



## Support



If you encounter any issues or have questions, please file an issue on the GitHub [issue tracker](https://github.com/yourusername/floom-cli/issues).
# Go Boilerplate

This Go boilerplate provides a simple and organized setup to kickstart your projects quickly. With a ready-to-go configuration, it comes packed with best practices, such as:

- Modular architecture
- Custom error handling
- Custom logging
- JWT authentication

## Installation

### Clone this project

```bash
gh repo clone bernardinorafael/go-boilerplate
```

### Rename the folder as you wish

```bash
mv go-boilerplate <new-name-here>
```

### Build Docker image

```bash
make docker-build
```

### Start containers

```bash
docker compose up -d
```

### Run migrations

```bash
make migrate-up
```

### Real-time logs

```bash
make air
```

## Authors

- [@bernardinodev](https://x.com/bernardinodev)

## Contribution

If you wish to contribute to this project, please follow the contribution guidelines and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

bash
#!/bin/bash
cd "$(dirname "$0")"
go build -o order-system main.go
echo "Build successful: order-system"
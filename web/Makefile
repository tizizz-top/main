# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=tizizz-web
WX_VERIFY=MP_verify_VmrbuqOYQnB6tzAS.txt
# Build the project
all: build

build:
	# build frontend
	# cd ui && pnpm build
	# cd ..
	cp -f ui/$(WX_VERIFY) ui/out/$(WX_VERIFY)
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Format the code
fmt:
	$(GOCMD) fmt ./...

proxy:
	# 将端口 8080 的流量转发到 192.168.1.100:80
	iptables -t nat -A PREROUTING -p tcp --dport 8080 -j DNAT --to-destination 192.168.1.100:80
	# 允许流量从本地网络转发
	iptables -t nat -A POSTROUTING -j MASQUERADE

rmproxy:
	iptables -t nat -L -n -v
	# 删除 PREROUTING 规则
	iptables -t nat -D PREROUTING -p tcp --dport 8080 -j DNAT --to-destination 192.168.1.100:80
	# 删除 POSTROUTING 规则
	iptables -t nat -D POSTROUTING -j MASQUERADE	

.PHONY: all build clean test deps fmt
# Cloud Infrastructure as Code

This repository demonstrates deploying OpenStack infrastructure in three ways: manually with DevStack, understanding the OpenStack architecture, and automating deployment with Terraform.

---

## 1. DevStack Installation (local.conf)

The **local.conf** file is used to automate OpenStack installation using DevStack, a development setup script.

### Setup

1. Copy the template file:
```bash
cp local.conf.example local.conf
```

2. Edit `local.conf` and set your own passwords:
```bash
nano local.conf
```

3. Run DevStack with this config file to deploy OpenStack.

### Configuration

The `local.conf` file contains:
- **Passwords** - Credentials for OpenStack services (keep these private!)
- **Network** - Automatically detects your host IP
- **Logging** - Configures log file location and retention
- **Reclone** - Set to `no` to skip re-cloning during updates

⚠️ **Important:** `local.conf` is in `.gitignore` - never commit actual passwords to version control.

---

## 2. OpenStack Architecture

### VM Creation Flow

![VM Creation Flow](./first_pic.png)

This diagram shows how a virtual machine is created in OpenStack:
1. **Horizon** - Web dashboard receives the request
2. **Keystone** - Authenticates the user
3. **Nova** - Processes the compute request
4. **Hypervisor** - Executes the VM creation

### Service Architecture

![Service Database Architecture](./second_pic.png)

OpenStack follows a modular architecture where:
- Each service (Keystone, Nova, Neutron, etc.) has its own database
- Each service manages its own state independently
- Services communicate via APIs and message queues
- This ensures fault isolation and scalability

---

## 3. Infrastructure as Code with Terraform

Instead of manually configuring OpenStack with DevStack, we use **Terraform** to define and deploy infrastructure automatically.

### What is main.tf?

The **main.tf** file contains:
- **Provider configuration** - Connects to OpenStack
- **Data sources** - Looks up images and flavors dynamically
- **Variables** - Customizable parameters for deployment
- **Resources** - Defines the compute instance to create

### Key Components

**Variables:**
- `instance_name` - Name of the VM (default: "basic")
- `network_name` - Network to attach to (default: "shared")
- `security_groups` - Firewall rules (default: ["default"])

**Resources:**
- Creates a Cirros 0.6.3 compute instance
- Deploys on m1.nano flavor
- Attaches to a specified network
- Applies security groups

### Usage

**Initialize Terraform:**
```bash
terraform init
```

**Plan the deployment:**
```bash
terraform plan
```

**Deploy:**
```bash
terraform apply
```

**Destroy resources:**
```bash
terraform destroy
```

**Deploy with custom variables:**
```bash
terraform apply -var="instance_name=my-vm" -var="security_groups=[\"default\",\"web\"]"
```

### Why Terraform Instead of Manual DevStack?

- **Reproducible** - Same infrastructure every time
- **Version controlled** - Track infrastructure changes in git
- **Scalable** - Deploy multiple instances easily
- **Documented** - Code serves as documentation
- **Idempotent** - Safe to run multiple times

---

## Prerequisites

- OpenStack environment (with credentials configured)
- [Terraform](https://www.terraform.io/downloads) >= 0.14.0

## Files

- **local.conf** - DevStack configuration for manual OpenStack installation
- **main.tf** - Terraform configuration for automated VM deployment
- **first_pic.png** - OpenStack VM creation flow diagram
- **second_pic.png** - OpenStack service architecture diagram

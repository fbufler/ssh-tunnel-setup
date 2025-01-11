# Easy Remote SSH Access

This guide explains how to set up an SSH tunnel using a dedicated VM with a static IP for easy remote access. Follow the steps in the given order.

## Step 1: Set Up the Server (VM with Static IP)

1. **Create a New User for the Tunnel:**

   ```bash
   sudo adduser tunneluser
   sudo usermod -aG sudo tunneluser
   ```

2. **Modify the SSH Configuration:**
   Open the SSH configuration file:

   ```bash
   sudo vim /etc/sshd_config
   ```

   Add the following lines:

   ```text
   GatewayPorts yes
   AllowTcpForwarding yes
   ```

3. **Restart the SSH Service:**

   ```bash
   sudo systemctl restart ssh
   ```

4. **Prepare the User's SSH Environment:**
   Switch to the new user and set up their SSH directory:
   ```bash
   su tunneluser
   mkdir ~/.ssh/
   chmod 700 ~/.ssh
   touch ~/.ssh/authorized_keys
   chmod 600 ~/.ssh/authorized_keys
   ```

## Step 2: Generate SSH Keys

### On the Target Machine

1. Generate an SSH key pair:
   ```bash
   ssh-keygen -t ed25519 -C "<some-name>"
   ```
2. Copy the public key:
   ```bash
   cat ~/.ssh/ed25519.pub
   ```

### On the Client Machine

1. Generate an SSH key pair:
   ```bash
   ssh-keygen -t ed25519 -C "<some-name>"
   ```
2. Copy the public key:
   ```bash
   cat ~/.ssh/ed25519.pub
   ```

## Step 3: Configure the Server to Accept Connections

1. **Add Public Keys to the Server:**
   Open the `authorized_keys` file for the `tunneluser`:
   ```bash
   vim ~/.ssh/authorized_keys
   ```
2. Paste the public keys:
   - The public key from the **Target Machine**.
   - The public key from the **Client Machine**.

## Step 4: Create the SSH Tunnel

### On the Target Machine

1. Establish the tunnel using the `-R` flag to reverse-forward the SSH port:
   ```bash
   ssh -R 2222:localhost:22 tunneluser@<vm-ip-or-dns>
   ```

### On the Client Machine

1. Connect to the target machine through the tunnel:
   ```bash
   ssh -p 2222 -i ~/.ssh/ed25519 <target-user>@<vm-ip-or-dns>
   ```
   - Enter the password if prompted.
   - You're done!

## Step 5: Make the Tunnel Persistent

### Automate Tunnel Creation

1. Create or modify the SSH config file (`~/.ssh/config`) on the **Target Machine**:
   ```text
   Host my_tunnel
       HostName <server_ip>
       User tunneluser
       IdentityFile ~/.ssh/id_ed25519
       RemoteForward 2222 localhost:22
   ```
2. Start the tunnel with:
   ```bash
   ssh my_tunnel
   ```

### Automate with Systemd

1. **Create a `systemd` Service on the Target Machine:**
   ```bash
   sudo nano /etc/systemd/system/ssh-tunnel.service
   ```
2. Add the following configuration:

   ```text
   [Unit]
   Description=SSH Tunnel
   After=network.target

   [Service]
   ExecStart=/usr/bin/ssh -N -R 2222:localhost:22 tunneluser@<vm-ip-or-dns>
   Restart=always
   User=<your_user>

   [Install]
   WantedBy=multi-user.target
   ```

3. **Enable and Start the Service:**
   ```bash
   sudo systemctl enable ssh-tunnel
   sudo systemctl start ssh-tunnel
   ```

### Monitor Tunnel Uptime with Cron

1. **Create a Monitoring Script:**

   ```bash
   nano /path/to/tunnel-check.sh
   ```

   Add the following:

   ```bash
   #!/bin/bash
   if ! nc -z localhost 2222; then
       systemctl restart ssh-tunnel
   fi
   ```

2. **Schedule the Script Using Cron:**
   ```bash
   crontab -e
   ```
   Add the following line:
   ```text
   */5 * * * * /path/to/tunnel-check.sh
   ```

## Step 6: Additional Measures

1. **Use ProxyCommand for Chained Tunnels (Optional):**
   If the server is not directly accessible, configure a jump host in the SSH config:

   ```text
   Host jump_host
       HostName <jump_host_ip>
       User <jump_user>

   Host target
       ProxyCommand ssh -W %h:%p jump_host
   ```

2. **Rotate Keys Regularly:**
   Periodically generate new SSH key pairs and update the configuration.

3. **Keep Software Updated:**
   Regularly update OpenSSH on the server, client, and target machines to apply security patches.

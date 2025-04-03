# Developer workflow:

a. Create a feature/bug branch from dev:

```
git checkout dev
git pull origin dev
git checkout -b feature/your-feature-name
```

b. Make changes and commit:

```
git add .
git commit -m "Describe your changes"
```

c. Push changes and create PR to dev:

```
git push origin feature/your-feature-name
```

- Create a Pull Request on GitHub from your feature branch to dev

d. After testing, create PR to staging:

- Create a Pull Request on GitHub from dev to staging

e. Create a new git tag:

```
git checkout staging
git pull origin staging
git tag -a v1.0.2-rc -m "Release version v1.0.2-rc"
git push origin v1.0.2-rc
```

# Deployment New Build/Release Process

To build the Paydoh Backend Project, follow these steps:

1. Clone the repository:

   ```
   git clone https://github.com/Grapesberry-Technologies-Pvt-Ltd/bank-api
   ```

2. Navigate to the project directory:

   ```
   cd bank-api

   <!-- install submodule (common service repository) -->
   1. git clone https://github.com/Grapesberry-Technologies-Pvt-Ltd/paydoh-commons.git
   2. git submodule init
   3. git submodule update --recursive
   ```

3. Check out the latest version:

   ```
   git checkout v1.0.1-rc
   ```

4. Build the project:

   ```
   make build
   ```

4.1. 

build-upload:
    # Step 1: SSH into the server and remove the old binary
    ssh -i "paydoh-key.pem" ubuntu@ec2-13-200-72-93.ap-south-1.compute.amazonaws.com 'cd bankapi && rm -rf bankapi'

makefile: 

    # Step 3: Provide permission to the new binary
  8.1.  ssh -i "paydoh-key.pem" ubuntu@ec2-13-200-72-93.ap-south-1.compute.amazonaws.com 'chmod +x bankapi'


5. Upload the binary to golang api server:

   ```
   make build-upload
   ```

6. Login to golang-api server:
   ```
   ssh -i "paydoh-key.pem" ubuntu@ec2-13-200-72-93.ap-south-1.compute.amazonaws.com
   ```
7. Go to bankapi directory:

   ```
   cd bankapi
   ```

8. Give the permission to the new uploaded binary:
   ```
   chmod +x bankapi
   ```
9. Stop the bankapi running demon:

   ```
   sudo systemctl stop bankapi
   ```

10. Start the bankapi demon now:

    ```
    sudo systemctl start bankapi
    ```

11. Check the status of the bankapi demon:

    ```
    sudo systemctl status bankapi
    ```

12. Check the logs of the bankapi demon:
    ```
    sudo journalctl -u bankapi
    ```

13. GIN_MODE=release ./bankapi #run in server
14. GIN_MODE=debug ./bankapi #run in local

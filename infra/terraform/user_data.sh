#!/bin/bash
# Update and install prerequisites
yum update -y
amazon-linux-extras install docker git -y
service docker start
usermod -a -G docker ec2-user

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Export environment variables
cat <<EOF >/etc/profile.d/luno_env.sh
export API_KEY_ID="${api_key_id}"
export API_KEY_SECRET="${api_key_secret}"
export DOCKERHUB_USERNAME="${dockerhub_username}"
export DOCKERHUB_TOKEN="${dockerhub_token}"
EOF

# Clone the repository and start services
git clone "${repo_url}" /opt/luno
cd /opt/luno
# Docker Hub login and deploy
docker login -u "${dockerhub_username}" -p "${dockerhub_token}"
docker-compose up -d --build

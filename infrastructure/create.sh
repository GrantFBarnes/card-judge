#!/usr/bin/env bash
################################################################################
# Create instance of Card Judge in Digital Ocean
################################################################################

set -e # exit on any command error

cd "$(dirname "$0")"

APP_NAME="card-judge"
DROPLET_NAME="$APP_NAME-database"

################################################################################
# get latest database backup

cd backups

BACKUP_GPG_FILE=$(ls *.gpg | tail -n 1)
BACKUP_GPG_PATH="$(pwd)/$BACKUP_GPG_FILE"
if [ ! -f "$BACKUP_GPG_PATH" ]; then
	echo "File not found: $BACKUP_GPG_PATH"
	exit 1
fi

BACKUP_SQL_PATH="${BACKUP_GPG_PATH::-4}"
rm -f $BACKUP_SQL_PATH
gpg -d --output $BACKUP_SQL_PATH $BACKUP_GPG_PATH

cd ..

################################################################################
# check expected files

if [ ! -f "$BACKUP_SQL_PATH" ]; then
	echo "File not found: $BACKUP_SQL_PATH"
	exit 1
fi

SETUP_SCRIPT_PATH="$(pwd)/templates/setup.sh"
if [ ! -f "$SETUP_SCRIPT_PATH" ]; then
	echo "File not found: $SETUP_SCRIPT_PATH"
	exit 1
fi

APP_SPEC_PATH="$(pwd)/templates/spec.yaml"
if [ ! -f "$APP_SPEC_PATH" ]; then
	echo "File not found: $APP_SPEC_PATH"
	exit 1
fi

################################################################################
# check environment variables

if [[ -z "$CARD_JUDGE_SQL_USER" ]]; then
	echo "Environment variable not found: CARD_JUDGE_SQL_USER"
	exit 1
fi

if [[ -z "$CARD_JUDGE_SQL_PASSWORD" ]]; then
	echo "Environment variable not found: CARD_JUDGE_SQL_PASSWORD"
	exit 1
fi

if [[ -z "$CARD_JUDGE_JWT_SECRET" ]]; then
	echo "Environment variable not found: CARD_JUDGE_JWT_SECRET"
	exit 1
fi

################################################################################
# sync fork with upstream if needed
REPO=$(grep 'repo:' "$APP_SPEC_PATH" | awk '{print $2}')
UPSTREAM_REPO=$(grep 'upstream_repo:' "$APP_SPEC_PATH" | awk '{print $2}')
if [[ "$REPO" != "$UPSTREAM_REPO" ]]; then
    echo "Syncing fork with upstream $UPSTREAM_REPO..."
    # Add upstream remote if it doesn't exist
    if ! git remote | grep -q upstream; then
        git remote add upstream https://github.com/$UPSTREAM_REPO.git
    fi
    git fetch upstream
    # Push upstream/main directly to origin/main without modifying local main
    git push origin upstream/main:main
else
    echo "Upstream repo is the same as repo, skipping sync."
fi

################################################################################
# get ssh key

echo "Which of the following SSH Keys should have access to the database droplet?"
doctl compute ssh-key list --format=Name --no-header
read -p "SSH Key Name: " SSH_KEY_NAME
if [[ -z "$SSH_KEY_NAME" ]]; then
	echo "SSH Key Name not provided"
	exit 1
fi

SSH_KEY_ID=$(doctl compute ssh-key list --format=ID,Name --no-header | grep $SSH_KEY_NAME | cut -d ' ' -f 1)
if [[ -z "$SSH_KEY_ID" ]]; then
	echo "SSH Key ID not found"
	exit 1
fi

################################################################################
# create droplet

echo "----------------------------------------"
echo "Creating Droplet..."

if doctl compute droplet list --format=Name --no-header | grep -q $DROPLET_NAME; then
	echo "Droplet already exists"
	exit 1
fi

sed -i -e 's/REPLACE_CARD_JUDGE_SQL_USER/'"$CARD_JUDGE_SQL_USER"'/g' "$SETUP_SCRIPT_PATH"
sed -i -e 's/REPLACE_CARD_JUDGE_SQL_PASSWORD/'"$CARD_JUDGE_SQL_PASSWORD"'/g' "$SETUP_SCRIPT_PATH"

DROPLET_IP=$(
	doctl compute droplet create "$DROPLET_NAME" \
		--ssh-keys=$SSH_KEY_ID \
		--region=nyc3 \
		--image=centos-stream-9-x64 \
		--size=s-1vcpu-1gb-amd \
		--user-data-file="$SETUP_SCRIPT_PATH" \
		--format=PublicIPv4 \
		--no-header \
		--wait
)

git checkout -- "$SETUP_SCRIPT_PATH"

if [[ -z "$DROPLET_IP" ]]; then
	sleep 10
	DROPLET_IP=$(doctl compute droplet list --format=PublicIPv4,Name --no-header | grep $DROPLET_NAME | cut -d ' ' -f 1)
	if [[ -z "$DROPLET_IP" ]]; then
		echo "Droplet IP not found"
		exit 1
	fi
fi

echo "Droplet Created"
echo "Droplet IP: $DROPLET_IP"

echo "Waiting 15 minutes for droplet to finish setup..."
sleep 15m

################################################################################
# restore database from backup

echo "----------------------------------------"
echo "Restoring Database..."

scp -o StrictHostKeyChecking=no "$BACKUP_SQL_PATH" root@$DROPLET_IP:/root/restore.sql >/dev/null 2>&1
ssh root@$DROPLET_IP 'mariadb CARD_JUDGE < /root/restore.sql'

echo "Database Restored"

################################################################################
# create app

echo "----------------------------------------"

if doctl apps list --format=Spec.Name --no-header | grep -q $APP_NAME; then
	echo "App already exists"
	exit 1
fi

sed -i -e 's/REPLACE_APP_NAME/'"$APP_NAME"'/g' "$APP_SPEC_PATH"
sed -i -e 's/REPLACE_CARD_JUDGE_SQL_HOST/'"$DROPLET_IP"'/g' "$APP_SPEC_PATH"
sed -i -e 's/REPLACE_CARD_JUDGE_SQL_USER/'"$CARD_JUDGE_SQL_USER"'/g' "$APP_SPEC_PATH"
sed -i -e 's/REPLACE_CARD_JUDGE_SQL_PASSWORD/'"$CARD_JUDGE_SQL_PASSWORD"'/g' "$APP_SPEC_PATH"
sed -i -e 's/REPLACE_CARD_JUDGE_JWT_SECRET/'"$CARD_JUDGE_JWT_SECRET"'/g' "$APP_SPEC_PATH"
sed -i -e 's/REPLACE_CARD_JUDGE_GIT_REPO/'"$CARD_JUDGE_GIT_REPO"'/g' "$APP_SPEC_PATH"

APP_URL=$(
	doctl apps create \
		--spec="$APP_SPEC_PATH" \
		--format=DefaultIngress \
		--no-header \
		--wait
)

git checkout -- "$APP_SPEC_PATH"

if [[ -z "$APP_URL" ]]; then
	sleep 10
	APP_URL=$(doctl apps list --format=DefaultIngress,Spec.Name --no-header | grep $APP_NAME | cut -d ' ' -f 1)
	if [[ -z "$APP_URL" ]]; then
		echo "App URL not found"
		exit 1
	fi
fi

echo "App URL: $APP_URL"

################################################################################

exit 0

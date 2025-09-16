#!/bin/bash

# Setup cron jobs for YouTube Analytics batch processing

set -e

echo "Setting up cron jobs for YouTube Analytics batch processing..."

# Check if running as root or with sudo
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root or with sudo"
   exit 1
fi

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Create log directory
LOG_DIR="/var/log/youtube-analytics"
mkdir -p "$LOG_DIR"
chown $SUDO_USER:$SUDO_USER "$LOG_DIR"

# Create cron wrapper script
cat > /usr/local/bin/youtube-analytics-batch <<EOF
#!/bin/bash
# Wrapper script for YouTube Analytics batch jobs

# Load environment variables
if [ -f /etc/youtube-analytics/env ]; then
    set -a
    source /etc/youtube-analytics/env
    set +a
fi

# Set working directory
cd $PROJECT_DIR

# Execute the batch command
exec "\$@" >> "$LOG_DIR/\$(basename \$1).log" 2>&1
EOF

chmod +x /usr/local/bin/youtube-analytics-batch

# Create environment file directory
mkdir -p /etc/youtube-analytics

# Copy environment variables (user must create /etc/youtube-analytics/env)
if [ ! -f /etc/youtube-analytics/env ]; then
    echo "Please create /etc/youtube-analytics/env with required environment variables:"
    echo "  DATABASE_URL=postgresql://..."
    echo "  YOUTUBE_API_KEY=..."
    echo "  PUBSUB_PROJECT_ID=..."
    echo "  CLOUDTASKS_PROJECT_ID=..."
    echo "  CLOUDTASKS_LOCATION=..."
    echo "  CLOUDTASKS_QUEUE_NAME=..."
    echo "  CLOUDTASKS_SERVICE_URL=..."
    echo "  WEBSUB_CALLBACK_URL=..."
fi

# Install cron jobs
CRON_FILE="/etc/cron.d/youtube-analytics"
cat > "$CRON_FILE" <<EOF
# YouTube Analytics Batch Processing Schedule
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

# Collect trending videos - Twice daily at 3:00 AM and 3:00 PM
0 3,15 * * * $SUDO_USER /usr/local/bin/youtube-analytics-batch $PROJECT_DIR/bin/batch-trending

# Schedule snapshots for recent videos - Every hour
0 * * * * $SUDO_USER /usr/local/bin/youtube-analytics-batch $PROJECT_DIR/bin/batch-schedule-snapshots -hours 2

# Generate rankings - Daily at 6:00 AM
0 6 * * * $SUDO_USER /usr/local/bin/youtube-analytics-batch $PROJECT_DIR/bin/batch-rankings

# Renew WebSub subscriptions - Daily at 1:00 AM
0 1 * * * $SUDO_USER /usr/local/bin/youtube-analytics-batch $PROJECT_DIR/bin/batch-websub-renewal -days 7

# Log rotation - Daily at midnight
0 0 * * * root find $LOG_DIR -name "*.log" -mtime +7 -delete
EOF

# Set proper permissions
chmod 0644 "$CRON_FILE"

# Build all batch binaries
echo "Building batch binaries..."
cd "$PROJECT_DIR"
make build-batch

# Restart cron service
echo "Restarting cron service..."
systemctl restart cron || service cron restart

echo "Cron jobs setup complete!"
echo ""
echo "Next steps:"
echo "1. Create /etc/youtube-analytics/env with required environment variables"
echo "2. Check logs in $LOG_DIR"
echo "3. Monitor cron execution: tail -f /var/log/syslog | grep CRON"
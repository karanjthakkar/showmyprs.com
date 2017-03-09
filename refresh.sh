# Crontab entry:
# 0,30 * * * * /home/ubuntu/showmyprs/refresh.sh >> /home/ubuntu/showmypr-cron-out.log 2>>/home/ubuntu/showmypr-cron-err.log
#
# This is a cronjob that runs every 30 minutes
# and updates the cache for the oldest 10 users.
echo "Running cronjob $(date)"

# Base url for the user profile endpoint
url="http://localhost:9001/user"

# Get a list of files inside the cache directory ordered by time (oldest first)
# Iterate through the top 10 of the filenames (which are usernames)
# Purge the cache and then call the user profile endpoint to 
# refill the cache.
for username in $(ls -hltr /home/ubuntu/showmyprs/.cache | awk '{print $9}' | head -10); do
  content="$(rm "/home/ubuntu/showmyprs/.cache/$username")"
  echo "Cleared cache for $username"
  content="$(curl -s "$url/$username")"
  echo "Refilling cache for $username"
done
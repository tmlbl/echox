# GET /
show_headers() {
  user_agent=$(echo "$headers" | grep -i user-agent | awk '{print $2}')
  echo "Your user agent is $user_agent"
}

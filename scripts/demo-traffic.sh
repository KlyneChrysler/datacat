#!/bin/bash
# E2e traffic personas against the local edge-proxy for DURATION seconds.
# scraper: metronome cadence, unique paths, curl UA  -> expect abusive
# human:   jittery cadence, few repeated paths, browser UA -> expect human
DURATION=${1:-90}
PROXY=http://localhost:8080
BROWSER_UA="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15"

scraper() {
	local end=$((SECONDS + DURATION)) i=0
	while ((SECONDS < end)); do
		curl -s -o /dev/null -b "dc_session=sim-scraper" "$PROXY/catalog/item-$i"
		i=$((i + 1))
		sleep 0.4
	done
}

human() {
	local end=$((SECONDS + DURATION))
	local paths=("/home" "/products" "/products" "/cart" "/home")
	while ((SECONDS < end)); do
		local path=${paths[RANDOM % ${#paths[@]}]}
		curl -s -o /dev/null -A "$BROWSER_UA" -b "dc_session=sim-human" "$PROXY$path"
		sleep "1.$((RANDOM % 9))"
		((RANDOM % 3 == 0)) && sleep $((RANDOM % 3))
	done
}

scraper &
human &
wait
echo "traffic done"

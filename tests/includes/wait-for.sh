# SHORT_TIMEOUT creates a consistent short timeout of the wait_for condition.
SHORT_TIMEOUT=5

# wait_for defines the ability to wait for a given condition to happen in a
# juju status output. The output is JSON, so everything that the API server
# knows about should be valid.
# The query argument is a jq query.
#
# ```
# wait_for <model name> <query>
# ```
wait_for() {
	local name query

	name=${1}
	query=${2}

	attempt=0
	# shellcheck disable=SC2046,SC2143
	until [[ "$(juju status --format=json 2>/dev/null | jq -S "${query}" | grep "${name}")" ]]; do
		echo "[+] (attempt ${attempt}) polling status for" "${query} => ${name}"
		juju status --relations 2>&1 | sed 's/^/    | /g'
		sleep "${SHORT_TIMEOUT}"
		attempt=$((attempt + 1))
	done

	if [[ ${attempt} -gt 0 ]]; then
		echo "[+] $(green 'Completed polling status for')" "$(green "${name}")"
		juju status --relations 2>&1 | sed 's/^/    | /g'
		# Although juju reports as an idle condition, some charms require a
		# breathe period to ensure things have actually settled.
		sleep "${SHORT_TIMEOUT}"
	fi
}

idle_condition() {
	local name app_index unit_index

	name=${1}
	app_index=${2:-0}
	unit_index=${3:-0}

	path=".[\"$name\"] | .units | .[\"$name/$unit_index\"]"

	echo ".applications | select(($path | .[\"juju-status\"] | .current == \"idle\") and ($path | .[\"workload-status\"] | .current != \"error\")) | keys[$app_index]"
}

idle_subordinate_condition() {
	local name parent unit_index

	name=${1}
	parent=${2}
	unit_index=${3:-0}

	path=".[\"$parent\"] | .units | .[] | .subordinates | .[\"$name/$unit_index\"]"

	# Print the *subordinate* name if it has an idle status in parent application
	echo ".applications | select(($path | .[\"juju-status\"] | .current == \"idle\") and ($path | .[\"workload-status\"] | .current != \"error\")) | \"$name\""
}

active_condition() {
	local name app_index unit_index

	name=${1}
	app_index=${2:-0}
	unit_index=${3:-0}

	echo ".applications | select(.[\"$name\"] | .[\"application-status\"] | .current == \"active\") | keys[$app_index]"
}

# workload_status gets the workload-status object for the unit - use
# .current or .message to select the actual field you need.
workload_status() {
	local app unit

	app=$1
	unit=$2

	echo ".applications[\"$app\"].units[\"$app/$unit\"][\"workload-status\"]"
}

# agent_status gets the juju-status object for the unit - use
# .current or .message to select the actual field you need.
agent_status() {
	local app unit

	app=$1
	unit=$2

	echo ".applications[\"$app\"].units[\"$app/$unit\"][\"juju-status\"]"
}

# charm_rev gets the current juju-status object for the application and uses it
# to find the application charm-rev.
charm_rev() {
	local app rev

	app=$1
	rev=${2:-0}

	echo ".applications | select(.[\"$app\"] | .[\"charm-rev\"] == $rev)"
}

# charm_channel gets the current juju-status object for the application and uses it
# to find the application charm-channel.
charm_channel() {
	local app channel

	app=$1
	channel=$2

	echo ".applications | select(.[\"$app\"] | .[\"charm-channel\"] == \"$channel\")"
}

# wait_for_machine_agent_status blocks until the machine agent for the specified
# machine instance ID reports the requested status.
#
# ```
# wait_for_machine_agent_status <instance-id> <status>
#
# example:
# wait_for_machine_agent_status "i-1234" "started"
# ```
wait_for_machine_agent_status() {
	local inst_id status

	inst_id=${1}
	status=${2}

	attempt=0
	# shellcheck disable=SC2046,SC2143
	until [ $(juju show-machine --format json | jq -r ".[\"machines\"] | .[\"${inst_id}\"] | .[\"juju-status\"] | .[\"current\"]" | grep "${status}") ]; do
		echo "[+] (attempt ${attempt}) polling machines"
		juju machines | grep "$inst_id" 2>&1 | sed 's/^/    | /g'
		sleep "${SHORT_TIMEOUT}"
		attempt=$((attempt + 1))
	done

	if [[ ${attempt} -gt 0 ]]; then
		echo "[+] $(green 'Completed polling machines')"
		juju machines | grep "$inst_id" 2>&1 | sed 's/^/    | /g'
		sleep "${SHORT_TIMEOUT}"
	fi
}

# wait_for_machine_netif_count blocks until the number of detected network
# interfaces for the requested machine instance ID becomes equal to the desired
# value.
#
# ```
# wait_for_machine_netif_count <instance-id> <count>
#
# example:
# wait_for_machine_netif_count "i-1234" "42"
# ```
wait_for_machine_netif_count() {
	local inst_id count

	inst_id=${1}
	count=${2}

	attempt=0
	# shellcheck disable=SC2046,SC2143
	until [ $(juju show-machine --format json | jq -r ".[\"machines\"] | .[\"${inst_id}\"] | .[\"network-interfaces\"] | length" | grep "${count}") ]; do
		# shellcheck disable=SC2046,SC2143
		echo "[+] (attempt ${attempt}) network interface count for instance ${inst_id} = "$(juju show-machine --format json | jq -r ".[\"machines\"] | .[\"${inst_id}\"] | .[\"network-interfaces\"] | length")
		sleep "${SHORT_TIMEOUT}"
		attempt=$((attempt + 1))
	done

}

# wait_for_subordinate_count blocks until the number of subordinates
# to the desired unit becomes equal to the desired value.
#
# ```
# wait_for_subordinate_count <application name> <principal unit num> <count>
#
# example:
# wait_for_subordinate_count mysql 0 3
# ```
wait_for_subordinate_count() {
	local name unit_index count

	name=${1}
	unit_index=${2:-0}
	count=${3:-0}

	attempt=0
	# shellcheck disable=SC2046,SC2143
	until [ $(juju status --format json | jq -r ".applications | .[\"${name}\"] | .units | .[\"${name}/${unit_index}\"] | .subordinates | length" | grep "${count}") ]; do
		# shellcheck disable=SC2046,SC2143
		echo "[+] (attempt ${attempt}) subordinate count for unit ${name}/${unit_index} = "$(juju status --format json | jq -r ".applications | .[\"${name}\"] | .units | .[\"${name}/${unit_index}\"] | .subordinates  | length")
		sleep "${SHORT_TIMEOUT}"
		attempt=$((attempt + 1))
	done

	if [[ ${attempt} -gt 0 ]]; then
		echo "[+] $(green 'Completed polling status')"
		juju status 2>&1 | sed 's/^/    | /g'
		sleep "${SHORT_TIMEOUT}"
	fi
}

# wait_for_model blocks until a model appears
# interfaces for the requested machine instance ID becomes equal to the desired
# value.
#
# ```
# wait_for_model <name>
#
# example:
# wait_for_model "default"
# ```
wait_for_model() {
	local name

	name=${1}

	attempt=0
	# shellcheck disable=SC2046,SC2143
	until [ $(juju models --format=json | jq -r ".models | .[] | select(.[\"short-name\"] == \"${name}\") | .[\"short-name\"]" | grep "${name}") ]; do
		echo "[+] (attempt ${attempt}) polling models"
		juju models | sed 's/^/    | /g'
		sleep "${SHORT_TIMEOUT}"
		attempt=$((attempt + 1))
	done

	if [[ ${attempt} -gt 0 ]]; then
		echo "[+] $(green 'Completed polling models')"
		juju models | sed 's/^/    | /g'
		sleep "${SHORT_TIMEOUT}"
	fi
}

# wait_for_systemd_service_files_to_appear blocks until the systemd service
# file for a unit is written to disk.
#
# ```
# wait_for_systemd_service_files_to_appear <unit_name>
#
# example:
# wait_for_systemd_service_files_to_appear "ubuntu-lite/0"
# ```
wait_for_systemd_service_files_to_appear() {
	local unit

	unit=${1}
	# shellcheck disable=SC2086
	svc_file_path="/etc/systemd/system/jujud-unit-$(echo -n ${1} | tr '/' '-').service"

	attempt=0
	# shellcheck disable=SC2046,SC2143
	while [ "$attempt" != "3" ]; do
		echo "[+] (attempt ${attempt}) waiting for the systemd unit files for ${unit} to appear"

		svc_present=$(juju ssh "${unit}" "ls ${svc_file_path} 2>/dev/null || echo -n 'missing'")
		if [[ ${svc_present} != "missing" ]]; then
			echo "[+] systemd unit files for ${unit} are now available"
			return
		fi

		sleep "${SHORT_TIMEOUT}"
		attempt=$((attempt + 1))
	done

	# shellcheck disable=SC2046
	echo $(red "Timed out waiting for the systemd unit files for ${unit} to appear")
	exit 1
}

run_refresh_switch_cs_to_ch() {
	# Test juju refresh from a charm store charm to a charm hub charm
	echo

	model_name="test-refresh-switch-ch"
	file="${TEST_DIR}/${model_name}.log"

	ensure "${model_name}" "${file}"

	juju deploy cs:ubuntu-19
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	OUT=$(juju refresh ubuntu --switch ch:ubuntu 2>&1 || true)
	if echo "${OUT}" | grep -E -vq "Added charm-hub charm"; then
		# shellcheck disable=SC2046
		echo $(red "failed refreshing charm: ${OUT}")
		exit 5
	fi
	# shellcheck disable=SC2059
	printf "${OUT}\n"

	# format: WARNING Charm store charms, those with cs: before the charm name, will not be supported in juju 3.1. Migration of this model to a juju 3.1 controller will be prohibited. Use the --switch flag to refresh to a non charm store version of the charm. Added charm-store charm "ubuntu", revision 21 in channel stable, to the model
	revision=$(echo "${OUT}" | awk 'BEGIN{FS=","} {print $4}' | awk 'BEGIN{FS=" "} {print $2}')

	wait_for "ubuntu" "$(charm_rev "ubuntu" "${revision}")"
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	destroy_model "${model_name}"
}

run_refresh_switch_cs_to_ch_no_new_revision() {
	# Test juju refresh from a charm store charm to a charm hub charm, with no new
	# charm revision.
	echo

	model_name="test-refresh-switch-ch"
	file="${TEST_DIR}/${model_name}.log"

	ensure "${model_name}" "${file}"

	OUT=$(juju deploy cs:ubuntu2 >&1 || true)
	# shellcheck disable=SC2059
	printf "${OUT}\n"
	# format: WARNING Charm store charms, those with cs: before the charm name, will not be supported in juju 3.1. Migration of this model to a juju 3.1 controller will be prohibited. Use the --switch flag to refresh to a non charm store version of the charm. Added charm-store charm "ubuntu", revision 21 in channel stable, to the model
	cs_revision=$(echo "${OUT}" | awk 'BEGIN{FS=","} {print $4}' | awk 'BEGIN{FS=" "} {print $2}')

	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	OUT=$(juju refresh ubuntu --switch ch:ubuntu 2>&1 || true)
	if echo "${OUT}" | grep -E -vq "Added charm-hub charm"; then
		# shellcheck disable=SC2046
		echo $(red "failed refreshing charm: ${OUT}")
		exit 5
	fi
	# shellcheck disable=SC2059
	printf "${OUT}\n"

	wait_for "ubuntu" "$(charm_rev "ubuntu" "${cs_revision}")"
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	destroy_model "${model_name}"
}

run_refresh_switch_cs_to_ch_channel() {
	# Test juju refresh from a charm store charm to a charm hub charm with a specific channel
	echo

	model_name="test-refresh-switch-ch-channel"
	file="${TEST_DIR}/${model_name}.log"

	ensure "${model_name}" "${file}"

	juju deploy cs:ubuntu-19
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	OUT=$(juju refresh ubuntu --switch ch:ubuntu --channel edge 2>&1 || true)
	if echo "${OUT}" | grep -E -vq "in channel edge"; then
		# shellcheck disable=SC2046
		echo $(red "failed refreshing charm: ${OUT}")
		exit 5
	fi
	# shellcheck disable=SC2059
	printf "${OUT}\n"

	# format: WARNING Charm store charms, those with cs: before the charm name, will not be supported in juju 3.1. Migration of this model to a juju 3.1 controller will be prohibited. Use the --switch flag to refresh to a non charm store version of the charm. Added charm-store charm "ubuntu", revision 21 in channel stable, to the model
	revision=$(echo "${OUT}" | awk 'BEGIN{FS=","} {print $4}' | awk 'BEGIN{FS=" "} {print $2}')

	wait_for "ubuntu" "$(charm_rev "ubuntu" "${revision}")"
	wait_for "ubuntu" "$(charm_channel "ubuntu" "edge")"
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	destroy_model "${model_name}"
}

run_refresh_switch_local_to_ch_channel() {
	# Test juju refresh from a local charm to a charm hub charm with a specific channel
	echo

	model_name="test-refresh-local-switch-ch"
	file="${TEST_DIR}/${model_name}.log"
	charm_name="${TEST_DIR}/ubuntu.charm"

	ensure "${model_name}" "${file}"

	juju download ubuntu --no-progress - >"${charm_name}"
	juju deploy "${charm_name}" ubuntu
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	OUT=$(juju refresh ubuntu --switch ch:ubuntu --channel edge 2>&1 || true)
	if echo "${OUT}" | grep -E -vq "Added charm-hub charm"; then
		# shellcheck disable=SC2046
		echo $(red "failed refreshing charm: ${OUT}")
		exit 5
	fi
	# shellcheck disable=SC2059
	printf "${OUT}\n"

	# format: WARNING Charm store charms, those with cs: before the charm name, will not be supported in juju 3.1. Migration of this model to a juju 3.1 controller will be prohibited. Use the --switch flag to refresh to a non charm store version of the charm. Added charm-store charm "ubuntu", revision 21 in channel stable, to the model
	revision=$(echo "${OUT}" | awk 'BEGIN{FS=","} {print $4}' | awk 'BEGIN{FS=" "} {print $2}')

	wait_for "ubuntu" "$(charm_rev "ubuntu" "${revision}")"
	wait_for "ubuntu" "$(charm_channel "ubuntu" "edge")"
	wait_for "ubuntu" "$(idle_condition "ubuntu")"

	destroy_model "${model_name}"
}

test_switch() {
	if [ "$(skip 'test_switch')" ]; then
		echo "==> TEST SKIPPED: refresh switch"
		return
	fi

	(
		set_verbosity

		cd .. || exit

		run "run_refresh_switch_cs_to_ch"
		run "run_refresh_switch_cs_to_ch_no_new_revision"
		run "run_refresh_switch_cs_to_ch_channel"
		run "run_refresh_switch_local_to_ch_channel"
	)
}

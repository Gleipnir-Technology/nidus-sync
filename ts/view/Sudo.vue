<template>
	<div class="container mt-4">
		<h2 id="comms-testing">
			<i class="bi bi-broadcast-pin"></i> Communications Testing
		</h2>
		<div class="row">
			<div class="col-lg-6">
				<SMSTest />
			</div>
			<div class="col-lg-6">
				<MMSTest />
			</div>
		</div>
		<div class="row">
			<div class="col-lg-6">
				<RCSTest />
			</div>
			<div class="col-lg-6">
				<EmailTest />
			</div>
		</div>
		<div class="row">
			<div class="col-lg-6">
				<SSETest />
			</div>
			<div class="col-lg-6">
				<PushNotificationTest />
			</div>
		</div>
		<div class="row">
			<UserImpersonation />
		</div>
	</div>
</template>

<script setup lang="ts">
import UserSelector from "@/components/UserSelector.vue";
import SMSTest from "@/components/sudo/SMSTest.vue";
import MMSTest from "@/components/sudo/MMSTest.vue";
import RCSTest from "@/components/sudo/RCSTest.vue";
import EmailTest from "@/components/sudo/EmailTest.vue";
import SSETest from "@/components/sudo/SSETest.vue";
import PushNotificationTest from "@/components/sudo/PushNotificationTest.vue";
import UserImpersonation from "@/components/sudo/UserImpersonation.vue";
</script>

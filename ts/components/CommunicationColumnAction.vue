<style scoped>
.actions-panel {
	height: 100%;
	overflow-y: auto;
}
</style>
<template>
	<div class="card shadow-sm h-100">
		<div class="card-header bg-light pane-header">Actions</div>
		<div class="card-body scroll-pane">
			<div v-if="loading" class="loading">Loading...</div>
			<div v-else>
				<div
					v-if="!selectedCommunication"
					class="h-100 d-flex flex-column align-items-center justify-content-center text-muted p-3"
				>
					<i class="bi bi-gear fs-1"></i>
					<p class="mt-2 text-center">
						Actions will appear here when a report is selected
					</p>
				</div>

				<div
					v-if="selectedCommunication"
					class="actions-panel d-flex flex-column"
				>
					<div class="p-3 flex-grow-1">
						<!-- Create Signal -->
						<div class="d-grid mb-3">
							<button class="btn btn-success btn-lg" @click="markSignal()">
								<i class="bi bi-plus-circle me-2"></i>Mark Signal
							</button>
							<small class="text-muted mt-1"
								>This report is useful signal</small
							>
						</div>

						<!-- Mark Invalid -->
						<div class="d-grid mb-3">
							<button class="btn btn-outline-danger" @click="markInvalid()">
								<i class="bi bi-x-circle me-2"></i>Mark Invalid
							</button>
							<small class="text-muted mt-1">This report isn't useful</small>
						</div>

						<hr />

						<!-- Message Reporter -->
						<div
							v-if="
								!(
									selectedCommunication?.public_report?.reporter.has_email ||
									selectedCommunication?.public_report?.reporter.has_phone
								)
							"
							class="mb-3"
						>
							<h6>
								<i class="bi bi-chat-dots"></i> No Reporter Communications
								Available
							</h6>
						</div>
						<div
							v-if="
								selectedCommunication?.public_report?.reporter.has_email ||
								selectedCommunication?.public_report?.reporter.has_phone
							"
							class="mb-3"
						>
							<h6><i class="bi bi-chat-dots"></i> Message Reporter</h6>
							<div class="mb-2">
								<label class="form-label small text-muted"
									>Quick Templates</label
								>
								<select
									class="form-select form-select-sm"
									@change="handleTemplateChange"
								>
									<option value="">Select a template...</option>
									<option value="received">Report Received</option>
									<option value="scheduled">Service Scheduled</option>
									<option value="completed">Service Completed</option>
									<option value="need_info">Need More Information</option>
								</select>
							</div>
							<textarea
								class="form-control mb-2"
								rows="5"
								v-model="messageText"
								placeholder="Type your message to the reporter..."
							></textarea>
							<div class="d-grid">
								<button
									class="btn btn-primary"
									@click="sendMessage()"
									:disabled="!messageText.trim()"
								>
									<i class="bi bi-send me-2"></i>Send Message
								</button>
							</div>
						</div>

						<hr />

						<!-- Report History -->
						<div>
							<h6><i class="bi bi-clock-history"></i> Activity Log</h6>
							<div class="small">
								<div
									v-for="(entry, index) in selectedCommunication?.public_report
										?.log || []"
									:key="index"
									class="border-start border-2 ps-2 mb-2"
								>
									<ListCardActivityLog :entry="entry" />
								</div>
								<div
									v-if="
										!selectedCommunication?.public_report?.log ||
										selectedCommunication?.public_report?.log.length === 0
									"
									class="text-muted"
								>
									No activity yet
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { Communication, User } from "@/type/api";
import ListCardActivityLog from "@/components/ListCardActivityLog.vue";
interface Emits {
	(e: "markSignal"): void;
	(e: "markInvalid"): void;
	(e: "sendMessage", message: string): void;
}
interface Props {
	loading: boolean;
	selectedCommunication: Communication | null;
}
const emit = defineEmits<Emits>();

const messageText = ref("");
const props = withDefaults(defineProps<Props>(), {});
function applyMessageTemplate(template: string) {
	const templates = {
		received: `Dear ${props.selectedCommunication?.public_report?.reporter.name || "Resident"},\n\nThank you for submitting your report to the Mosquito Control District. We have received your communication and it has been assigned to our team for review.\n\nWe will be in touch if we need any additional information.\n\nBest regards,\nMosquito Control District`,
		scheduled: `Dear ${props.selectedCommunication?.public_report?.reporter.name || "Resident"},\n\nGood news! Based on your report, we have scheduled a service visit to your area. Our technicians will be conducting mosquito control operations within the next 3-5 business days.\n\nNo action is required on your part.\n\nBest regards,\nMosquito Control District`,
		completed: `Dear ${props.selectedCommunication?.public_report?.reporter.name || "Resident"},\n\nWe wanted to let you know that our team has completed mosquito control operations in your area based on your recent report.\n\nIf you continue to experience issues, please don't hesitate to submit a new report.\n\nBest regards,\nMosquito Control District`,
		need_info: `Dear ${props.selectedCommunication?.public_report?.reporter.name || "Resident"},\n\nThank you for your recent report. In order to better assist you, we need some additional information:\n\n- [Specify what information is needed]\n\nPlease reply to this message with the requested details.\n\nBest regards,\nMosquito Control District`,
	};

	if (template in templates) {
		messageText.value = templates[template as keyof typeof templates];
	}
}
function handleTemplateChange(event: Event) {
	const target = event.target as HTMLSelectElement;
	applyMessageTemplate(target.value);
}
function markInvalid() {
	emit("markInvalid");
}
function markSignal() {
	emit("markSignal");
}
function sendMessage() {
	emit("sendMessage", messageText.value);
}
</script>

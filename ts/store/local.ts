import { defineStore } from "pinia";

export const useStoreLocal = defineStore("local", () => {
	function getSessionID(): string {
		let id = localStorage.getItem("session_id");
		if (id) {
			return id;
		}
		id = crypto.randomUUID();
		localStorage.setItem("session_id", id.toString());
		return id;
	}
	return {
		getSessionID,
	};
});

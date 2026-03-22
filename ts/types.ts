export interface Address {
	locality: string;
	number: string;
	street: string;
}

export interface PublicReport {
	created: string;
	type: string;
}

export interface Communication {
	created: string;
	id: string;
	public_report: PublicReport | null;
	type: string;
}

type Source = {
	source: string;
}
type ListData = {
	description: string;
	sources: [Source]
}

export function getList(token: string) {
	return fetch(`/api/list/${token}`)
		.then(res => res.json()).then<ListData>(json => json);
}

export function updateDescription(token: string, description: string) {
	return fetch(`/api/list/${token}?description=${description}`, {
		method: 'POST'
	});
}

export function addSource(token: string, source: string) {
	return fetch(`/api/list/${token}?source=${source}`, {
		method: 'POST'
	});
}

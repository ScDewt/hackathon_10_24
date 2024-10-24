import React from 'react';
import { Button } from 'antd';

export function CreateList() {
	const onClick = () => {
		const id = Math.random().toString(36).slice(2, 9);
		window.location.replace(`/id/${id}`);
	}
	return <Button type="primary" onClick={onClick}>Create list</Button>
}

import React from 'react';
import './App.css';
import { List } from './List';
import { CreateList } from './CreateList';

function App() {
	const currentUrl = window.location.href;
	const url = new URL(currentUrl);
	const pathElements = url.pathname.split('/');

	const token = pathElements[pathElements.length - 1];

	return (
		<div className="App">
			{ !token ? <CreateList /> : <List id={token} /> }
		</div>
	);
}

export default App;

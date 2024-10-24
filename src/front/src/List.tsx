import React, { useEffect, useRef, useState } from 'react';
import { Button, Form, Input, Space, List as UiList, TimePicker, Card } from 'antd';
import dayjs from 'dayjs';
import { addSource, getList, updateDescription } from './service';

type Props = {
	id: string
}
type InitValues = {
	description?: string;
	sources?: string[];
}

export function List({ id }: Props) {
	const [urls, setUrls] = useState<string[]>([]);
	const [time, setTime] = useState<string>('');
	const [description, setDescription] = useState<string>('');
	const [telegramUrl, setTelegramUrl] = useState<string>('');
	const initValues = useRef<InitValues>({});

	const addTelegramUrl = () => {
		setUrls([...urls, telegramUrl]);
		setTelegramUrl('');
	}

	useEffect(() => {
		getList(id).then(({ description, sources }) => {
			if (description) {
				initValues.current.description = description;
				setDescription(description);
			}
			if (Array.isArray(sources)) {
				sources.forEach(({ source }) => {
					if (!initValues.current.sources) {
						initValues.current.sources = [];
					}
					initValues.current.sources.push(source);
					setUrls([...urls, source])
				})
			}
		})
	}, []);



	const reset = () => {
		setTime('');
	}

	const send = () => {
		if (description && description !== initValues.current.description) {
			updateDescription(id, description)
		}
		const initSourcesSet = new Set(initValues.current.sources);
		urls.forEach((url) => {
			if (initSourcesSet.has(url)) {
				return;
			}
			addSource(id, url);
		})
		reset();
	}

	const format = 'HH:mm';

	return (
		<Form name="listform" style={{ maxWidth: 900 }}>
			<Card>
				<Space align="start">
					<Space direction="vertical">
						<UiList
							size="default"
							bordered
							dataSource={ urls }
							renderItem={ (item) => (
								<UiList.Item >
									{ item }
								</UiList.Item>
							) }
						/>

						<Form.Item name="note">
							<Space>
								<Input
									onChange={event => setTelegramUrl(event.target.value) }
									value={ telegramUrl }
									placeholder="Add Telegram link"
								/>
								<Button
									type="default"
									onClick={ addTelegramUrl }
								>
									Add
								</Button>
							</Space>
						</Form.Item>
					</Space>

					<Space direction="vertical">
						<div>
							<TimePicker
								defaultValue={dayjs('12:08', format)}
								format={format}
								onChange={(_, dateString) => setTime(Array.isArray(dateString) ? dateString[0] : dateString)}
							/>
						</div>
						<div>
							<Input.TextArea
								rows={4}
								placeholder="Description"
								onChange={event => setDescription(event.target.value)}
							/>
						</div>
					</Space>
				</Space>
				<div>
					<Button
						size="large"
						type="primary"
						onClick={ send }
					>
						Save
					</Button>
				</div>
			</Card>
		</Form>
	)
}

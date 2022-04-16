import { useEffect, useState } from 'react'
import axios from 'axios'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Loading from './states/Loading'
import Finished from './states/Finished'
import Warning from './states/Warning'

export default function Notification(props: any) {
	const [title, setTitle] = useState<string>(),
		[notificationClass, setNotificationClass] = useState('hidden'),
		[notificationStateClass, setNotificationStateClass] = useState<string>(),
		[state, setState] = useState<string>(),
		[action, setAction] = useState<any>(),
		[notify, setNotify] = useState(false)

	let getNotificationsInterval: any

	useEffect(() => {
		!getNotificationsInterval && (getNotificationsInterval = setInterval(() => getNotifications(), 2000))
	}, [])

	const getNotifications = () => {
		axios
			.request({
				method: 'GET',
				url: `http://4.2.0.225:5000/live/notify`,
				headers: {
					'Content-Type': 'application/json',
				},
			})
			.then((res) => {
				if (res.data.state !== 'inactive') {
					setTitle(res.data.title)
					setState(res.data.state)
					setAction(res.data.action)
					setNotify(true)
				} else {
					setNotify(false)
				}
			})
			.catch((error) => {
				console.error(error)
			})
	}

	useEffect(() => {
		if (notify) {
			setNotificationClass('notification')

			switch (state) {
				case 'inProgress': setNotificationStateClass('var(--irrigationBlue)'); break
				case 'finished': setNotificationStateClass('var(--green)'); break
				case 'physicalHelpRequired': setNotificationStateClass('var(--warningRed)'); break
			}
		} else {
			setNotificationClass('hidden')
		}
	}, [notify])

	return (
		<div className={notificationClass}>
			<Card className="card p-0-i drop-shadow-2xl">
				<CardContent className="p-0-i">
					<div className="flex-row p-2 pl-4" style={{ background: notificationStateClass }}>
						<span className="text-white title-2 font-semibold">{title}</span>
					</div>
					<div className="flex-row p-2">
						<div className="flex-col p-2 w-3/12">
							<div className="flex-row h-8">
								<span>
									{state === 'inProgress' && <Loading />}
									{state === 'finished' && <Finished />}
									{state === 'physicalHelpRequired' && <Warning />}
								</span>
							</div>
						</div>
						<div className="flex-col p-2 w-9/12">
							<div className="flex-row">
								<span className="title-2" style={{ color: notificationStateClass }}>
									{action}
								</span>
							</div>
						</div>
					</div>
				</CardContent>
			</Card>
		</div>
	)
}

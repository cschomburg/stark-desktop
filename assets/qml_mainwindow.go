package assets

const QmlMainWindow = `
import QtQuick 2.2
import QtQuick.Controls 1.2
import QtQuick.Layouts 1.0

ApplicationWindow {
	id: mainWindow
	title: "stark"
	minimumWidth: 400
	minimumHeight: 300
	width: 600
	height: 400

	signal publish(string text)

	ColumnLayout {
		anchors.fill: parent
		spacing: 4

		TextArea {
			id: messages
			Layout.fillWidth: true
			Layout.fillHeight: true
			text: ""
			textFormat: TextEdit.RichText
			readOnly: true
		}

		RowLayout {
			Layout.fillWidth: true

			TextField {
				id: replyInput
				Layout.fillWidth: true
				Layout.minimumHeight: 30
				text: ""
				placeholderText: "Talk to me ..."
				focus: true
			}

			Button {
				id: replySend
				text: "Send"
			}
		}
	}

	function setHistory(text) {
		messages.text = text
	}

	function publishReply(text) {
		var text = replyInput.text
		publish(text)
		replyInput.text = ""
	}

	Component.onCompleted: {
		replyInput.accepted.connect(publishReply)
		replySend.clicked.connect(publishReply)
	}
}
`

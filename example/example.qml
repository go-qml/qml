import QtQuick 2.0;
import QtGraphicalEffects 1.0;

Rectangle {
	width: 400
	height: 400
	gradient: Gradient {
		GradientStop { position: 0.0; color: "#3a2c32"; }
		GradientStop { position: 0.8; color: "#875864"; }
		GradientStop { position: 1.0; color: "#9b616c"; }
	}

	Text {
		text: message.text
		anchors.centerIn: parent
		color: "white"
		font.bold: true
		font.pointSize: 20
	}
}

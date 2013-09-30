import QtQuick 2.0

Item {
    width: 320; height: 200

    Component {
        id: itemDelegate
        Text {
		text: "I am color number: " + index
		color: colors.name(index)
	}
    }

    ListView {
        anchors.fill: parent
        model: colors.len()
        delegate: itemDelegate
    }
}

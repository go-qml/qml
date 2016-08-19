import QtQuick 2.0

Item {
    width: 320; height: 200

    Component.onCompleted: {
      console.log("colors:", colors);
    }

    ListView {
        width: 120;
        model: colors
        delegate: Text {
            text: "I am color number: " + index
            color: display
        }
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        anchors.horizontalCenter: parent.horizontalCenter
    }
}

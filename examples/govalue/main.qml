import QtQuick 2.0
Rectangle {
  width:600
  height:600
  Component.onCompleted: {
    var gt = gostruct.returnGoType()
    gt.useGoType(gt)
  }
}

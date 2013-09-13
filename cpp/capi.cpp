#include <QApplication>
#include <QQuickView>
#include <QtQml>
#include <QDebug>

#include "govalue.h"
#include "govaluetype.h"
#include "capi.h"

#include <QDebug>

void newGuiApplication()
{
    static char empty[1] = {0};
    static char *argv[] = {empty};
    static int argc = 1;
    new QGuiApplication(argc, argv);

    // The event should never die.
    qApp->setQuitOnLastWindowClosed(false);
}

void applicationExec()
{
    qApp->exec();
}

void applicationFlushAll()
{
    qApp->processEvents();
}

void *currentThread()
{
    return QThread::currentThread();
}

void *appThread()
{
    return QCoreApplication::instance()->thread();
}

QQmlEngine_ *newEngine(QObject_ *parent)
{
    return new QQmlEngine(reinterpret_cast<QObject *>(parent));
}

QQmlContext_ *engineRootContext(QQmlEngine_ *engine)
{
    return reinterpret_cast<QQmlEngine *>(engine)->rootContext();
}

void engineSetContextForObject(QQmlEngine_ *engine, QObject_ *object)
{
    QQmlEngine *qengine = reinterpret_cast<QQmlEngine *>(engine);
    QObject *qobject = reinterpret_cast<QObject *>(object);

    QQmlEngine::setContextForObject(qobject, qengine->rootContext());
}

void engineSetOwnershipCPP(QQmlEngine_ *engine, QObject_ *object)
{
    QQmlEngine *qengine = reinterpret_cast<QQmlEngine *>(engine);
    QObject *qobject = reinterpret_cast<QObject *>(object);

    qengine->setObjectOwnership(qobject, QQmlEngine::CppOwnership);
}

void engineSetOwnershipJS(QQmlEngine_ *engine, QObject_ *object)
{
    QQmlEngine *qengine = reinterpret_cast<QQmlEngine *>(engine);
    QObject *qobject = reinterpret_cast<QObject *>(object);

    qengine->setObjectOwnership(qobject, QQmlEngine::JavaScriptOwnership);
}

QQmlComponent_ *newComponent(QQmlEngine_ *engine, QObject_ *parent)
{
    QQmlEngine *qengine = reinterpret_cast<QQmlEngine *>(engine);
    //QObject *qparent = reinterpret_cast<QObject *>(parent);
    return new QQmlComponent(qengine);
}

void componentSetData(QQmlComponent_ *component, const char *data, int dataLen, const char *url, int urlLen)
{
    QByteArray qdata(data, dataLen);
    QByteArray qurl(url, urlLen);
    QString qsurl = QString::fromUtf8(qurl);
    reinterpret_cast<QQmlComponent *>(component)->setData(qdata, qsurl);
}

char *componentErrorString(QQmlComponent_ *component)
{
    QQmlComponent *qcomponent = reinterpret_cast<QQmlComponent *>(component);
    if (qcomponent->isReady()) {
        return NULL;
    }
    if (qcomponent->isError()) {
        QByteArray ba = qcomponent->errorString().toUtf8();
        return strdup(ba.constData());
    }
    return strdup("component is not ready (why!?)");
}

QObject_ *componentCreate(QQmlComponent_ *component, QQmlContext_ *context)
{
    QQmlComponent *qcomponent = reinterpret_cast<QQmlComponent *>(component);
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);

    QObject *instance = qcomponent->create(qcontext);
    return instance;
}

QQuickView_ *componentCreateView(QQmlComponent_ *component, QQmlContext_ *context)
{
    QQmlComponent *qcomponent = reinterpret_cast<QQmlComponent *>(component);
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);

    QObject *instance = qcomponent->create(qcontext);
    QQuickView *view = new QQuickView(qcontext->engine(), 0);
    view->setContent(qcomponent->url(), qcomponent, instance);
    view->setResizeMode(QQuickView::SizeRootObjectToView);
    return view;
}

void viewShow(QQuickView_ *view)
{
    reinterpret_cast<QQuickView *>(view)->show();
}

void viewHide(QQuickView_ *view)
{
    reinterpret_cast<QQuickView *>(view)->hide();
}

void viewConnectHidden(QQuickView_ *view)
{
    QQuickView *qview = reinterpret_cast<QQuickView *>(view);
    QObject::connect(qview, &QWindow::visibleChanged, [=](bool visible){
        if (!visible) {
            hookWindowHidden(view);
        }
    });
}

QObject_ *viewRootObject(QQuickView_ *view)
{
    QQuickView *qview = reinterpret_cast<QQuickView *>(view);
    return qview->rootObject();
}

void contextSetObject(QQmlContext_ *context, QObject_ *value)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    QObject *qvalue = reinterpret_cast<QObject *>(value);

    // Give qvalue an engine reference if it doesn't yet have one.
    if (!qmlEngine(qvalue)) {
        QQmlEngine::setContextForObject(qvalue, qcontext->engine()->rootContext());
    }

    qcontext->setContextObject(qvalue);
}

void contextSetProperty(QQmlContext_ *context, QString_ *name, DataValue *value)
{
    const QString *qname = reinterpret_cast<QString *>(name);
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);

    QVariant var;
    unpackDataValue(value, &var);

    // Give qvalue an engine reference if it doesn't yet have one .
    QObject *obj = var.value<QObject *>();
    if (obj && !qmlEngine(obj)) {
        QQmlEngine::setContextForObject(obj, qcontext);
    }

    qcontext->setContextProperty(*qname, var);
}

void contextGetProperty(QQmlContext_ *context, QString_ *name, DataValue *result)
{
    QQmlContext *qcontext = reinterpret_cast<QQmlContext *>(context);
    const QString *qname = reinterpret_cast<QString *>(name);

    QVariant var = qcontext->contextProperty(*qname);
    packDataValue(&var, result);
}

void delObject(QObject_ *object)
{
    delete reinterpret_cast<QObject *>(object);

    // That looks handy, but doesn't work well. Often objects will stay undeleted
    // for whatever reason and break the tests on the leak prevention.
    //reinterpret_cast<QObject *>(object)->deleteLater();
}

void objectGetProperty(QObject_ *object, const char *name, DataValue *value)
{
    QObject *qobject = reinterpret_cast<QObject *>(object);
    
    QVariant var = qobject->property(name);
    packDataValue(&var, value);
}

void objectSetProperty(QObject_ *object, const char *name, DataValue *value)
{
    QObject *qobject = reinterpret_cast<QObject *>(object);
    QVariant var;
    unpackDataValue(value, &var);

    // Give qvalue an engine reference if it doesn't yet have one.
    QObject *obj = var.value<QObject *>();
    if (obj && !qmlEngine(obj)) {
        QQmlContext *context = qmlContext(qobject);
        if (context) {
            QQmlEngine::setContextForObject(obj, context);
        }
    }

    qobject->setProperty(name, var);
}

void objectInvoke(QObject_ *object, const char *method, DataValue *resultdv, DataValue *paramsdv, int paramsLen)
{
    QObject *qobject = reinterpret_cast<QObject *>(object);

    QVariant result;
    QVariant param[MaximumParamCount-1];
    QGenericArgument arg[MaximumParamCount-1];
    for (int i = 0; i < paramsLen; i++) {
        unpackDataValue(&paramsdv[i], &param[i]);
        arg[i] = Q_ARG(QVariant, param[i]);
    }
    if (paramsLen > 10) {
        qFatal("fix the parameter dispatching");
    }
    QMetaObject::invokeMethod(qobject, method, Qt::DirectConnection, 
            Q_RETURN_ARG(QVariant, result), arg[0], arg[1], arg[2], arg[3], arg[4], arg[5], arg[6], arg[7], arg[8], arg[9]);
    packDataValue(&result, resultdv);
}

void objectFindChild(QObject_ *object, QString_ *name, DataValue *resultdv)
{
    QObject *qobject = reinterpret_cast<QObject *>(object);
    QString *qname = reinterpret_cast<QString *>(name);
    
    QVariant var;
    QObject *result = qobject->findChild<QObject *>(*qname);
    if (result) {
        var.setValue(result);
    }
    packDataValue(&var, resultdv);
}

void objectSetParent(QObject_ *object, QObject_ *parent)
{
    QObject *qobject = reinterpret_cast<QObject *>(object);
    QObject *qparent = reinterpret_cast<QObject *>(parent);

    qobject->setParent(qparent);
}

QString_ *newString(const char *data, int len)
{
    // This will copy data only once.
    QByteArray ba = QByteArray::fromRawData(data, len);
    return new QString(ba);
}

void delString(QString_ *s)
{
    delete reinterpret_cast<QString *>(s);
}

GoValue_ *newGoValue(GoAddr *addr, GoTypeInfo *typeInfo, QObject_ *parent)
{
    QObject *qparent = reinterpret_cast<QObject *>(parent);
    return new GoValue(addr, typeInfo, qparent);
}

void goValueActivate(GoValue_ *value, GoTypeInfo *typeInfo, int addrOffset)
{
    GoMemberInfo *fieldInfo = typeInfo->fields;
    for (int i = 0; i < typeInfo->fieldsLen; i++) {
        if (fieldInfo->addrOffset == addrOffset) {
            reinterpret_cast<GoValue *>(value)->activate(fieldInfo->metaIndex);
            return;
        }
        fieldInfo++;
    }

    // TODO Return an error; probably an unexported field.
}

template<int N>
void registerSingletonN(char *location, int major, int minor, char *name, GoTypeInfo *info, GoTypeSpec_ *spec) {
    GoValueType<N>::init(info, spec);
    qmlRegisterSingletonType< GoValueType<N> >(location, major, minor, name, [](QQmlEngine *qmlEngine, QJSEngine *jsEngine) -> QObject* {
        QObject *singleton = new GoValueType<N>();
        QQmlEngine::setContextForObject(singleton, qmlEngine->rootContext());
        return singleton;
    });
}

void registerSingleton(char *location, int major, int minor, char *name, GoTypeInfo *info, GoTypeSpec_ *spec)
{
    // TODO Must increment the number and use different types per call.
    registerSingletonN<1>(location, major, minor, name, info, spec);
}

void registerType(char *location, int major, int minor, char *name, GoTypeInfo *info, GoTypeSpec_ *spec)
{
    // TODO Must increment the number and use different types per call.
    GoValueType<1>::init(info, spec);
    qmlRegisterType< GoValueType<1> >(location, major, minor, name);
}

void unpackDataValue(DataValue *value, QVariant_ *var)
{
    QVariant *qvar = reinterpret_cast<QVariant *>(var);
    switch (value->dataType) {
    case DTString:
        *qvar = QString::fromUtf8(*(char **)value->data, value->len);
        break;
    case DTBool:
        *qvar = bool(*(char *)(value->data) != 0);
        break;
    case DTInt64:
        *qvar = *(qint64*)(value->data);
        break;
    case DTInt32:
        *qvar = *(qint32*)(value->data);
        break;
    case DTFloat64:
        *qvar = *(double*)(value->data);
        break;
    case DTFloat32:
        *qvar = *(float*)(value->data);
        break;
    case DTList:
        *qvar = **(QVariantList**)(value->data);
        delete *(QVariantList**)(value->data);
        break;
    case DTObject:
        qvar->setValue(*(QObject**)(value->data));
        break;
    case DTInvalid:
        // null would be more natural, but an invalid variant means
        // it has proper semantics when dealing with non-qml qt code.
        //qvar->setValue(QJSValue(QJSValue::NullValue));
        qvar->clear();
        break;
    default:
        qFatal("Unsupported data type: %d", value->dataType);
        break;
    }
}

void packDataValue(QVariant_ *var, DataValue *value)
{
    QVariant *qvar = reinterpret_cast<QVariant *>(var);

    // Some assumptions are made below regarding the size of types.
    // There's apparently no better way to handle this since that's
    // how the types with well defined sizes (qint64) are mapped to
    // meta-types (QMetaType::LongLong).
    switch (qvar->type()) {
    case QVariant::Invalid:
        value->dataType = DTInvalid;
        break;
    case QMetaType::QString:
        {
            value->dataType = DTString;
            QByteArray ba = qvar->toByteArray();
            *(char**)(value->data) = strdup(ba.constData());
            value->len = ba.size();
            break;
        }
    case QMetaType::Bool:
        value->dataType = DTBool;
        *(qint8*)(value->data) = (qint8)qvar->toInt();
        break;
    case QMetaType::LongLong:
        value->dataType = DTInt64;
        *(qint64*)(value->data) = qvar->toLongLong();
        break;
    case QMetaType::Int:
        value->dataType = DTInt32;
        *(qint32*)(value->data) = qvar->toInt();
        break;
    case QMetaType::Double:
        value->dataType = DTFloat64;
        *(double*)(value->data) = qvar->toDouble();
        break;
    case QMetaType::Float:
        value->dataType = DTFloat32;
        *(float*)(value->data) = qvar->toFloat();
        break;
    case QMetaType::QObjectStar:
        {
            QObject *qobject = qvar->value<QObject *>();
            GoValue *govalue = dynamic_cast<GoValue *>(qobject);
            if (govalue) {
                value->dataType = DTGoAddr;
                *(void **)(value->data) = govalue->addr();
            } else {
                value->dataType = DTObject;
                *(void **)(value->data) = qobject;
            }
            break;
        }
        // fallthrough
    default:
        qFatal("Unsupported variant type: %d", qvar->type());
        break;
    }
}

QVariantList_ *newVariantList(DataValue *list, int len)
{
    QVariantList *vlist = new QVariantList();
    vlist->reserve(len);
    for (int i = 0; i < len; i++) {
        QVariant var;
        unpackDataValue(&list[i], &var);
        vlist->append(var);
    }
    return vlist;
}

void internalLogHandler(QtMsgType severity, const QMessageLogContext &context, const QString &text)
{
    QByteArray textba = text.toUtf8();
    LogMessage message = {severity, textba.constData(), textba.size(), context.file, (int)strlen(context.file), context.line};
    hookLogHandler(&message);
}

void installLogHandler()
{
    qInstallMessageHandler(internalLogHandler);
}


#include <QGuiApplication>
#include <QtQuick/QQuickView>
#include <QtQml/QtQml>

#include <unistd.h>

extern "C" {

const char *particleQml =
"import QtQuick 2.0\n"
"import QtQuick.Particles 2.0\n"
"import QtGraphicalEffects 1.0;\n"
"\n"
"Rectangle {\n"
"	id: root\n"
"	color: \"black\"\n"
"\n"
"	width: 640\n"
"	height: 480\n"
"\n"
"	gradient: Gradient {\n"
"		GradientStop { position: 0.0; color: \"#3a2c32\"; }\n"
"		GradientStop { position: 0.8; color: \"#875864\"; }\n"
"		GradientStop { position: 1.0; color: \"#9b616c\"; }\n"
"	}\n"
"\n"
"	Text {\n"
"		//text: message.text\n"
"		text: \"Hello from Go!\"\n"
"\n"
"		x: parent.width/2 - width/2\n"
"		y: parent.height/2 - height/2\n"
"\n"
"		color: \"white\"\n"
"		font.bold: true\n"
"		font.pointSize: 20\n"
"\n"
"		MouseArea {\n"
"		    id: mouseArea\n"
"		    anchors.fill: parent\n"
"		    drag.target: parent\n"
"                    onReleased: {\n"
"			root.customEmit(parent.x + 15, parent.y + parent.height/2);\n"
"			root.customEmit(parent.x + parent.width / 2, parent.y + parent.height/2);\n"
"			root.customEmit(parent.x + parent.width - 15, parent.y + parent.height/2);\n"
"		    }\n"
"		}\n"
"	}\n"
"\n"
"	ParticleSystem {\n"
"		id: sys\n"
"	}\n"
"\n"
"	ImageParticle {\n"
"		system: sys\n"
"		source: \"particle.png\"\n"
"		color: \"white\"\n"
"		colorVariation: 1.0\n"
"		alpha: 0.1\n"
"	}\n"
"\n"
"	Component {\n"
"		id: emitterComp\n"
"		Emitter {\n"
"			id: container\n"
"			Emitter {\n"
"				id: emitMore\n"
"				system: sys\n"
"				emitRate: 128\n"
"				lifeSpan: 600\n"
"				size: 16\n"
"				endSize: 8\n"
"				velocity: AngleDirection {angleVariation:360; magnitude: 60}\n"
"			}\n"
"\n"
"			property int life: 2600\n"
"			property real targetX: 0\n"
"			property real targetY: 0\n"
"			function go() {\n"
"				xAnim.start();\n"
"				yAnim.start();\n"
"				container.enabled = true\n"
"			}\n"
"			system: sys\n"
"			emitRate: 32\n"
"			lifeSpan: 600\n"
"			size: 24\n"
"			endSize: 8\n"
"			NumberAnimation on x {\n"
"				id: xAnim;\n"
"				to: targetX\n"
"				duration: life\n"
"				running: false\n"
"			}\n"
"			NumberAnimation on y {\n"
"				id: yAnim;\n"
"				to: targetY\n"
"				duration: life\n"
"				running: false\n"
"			}\n"
"			Timer {\n"
"				interval: life\n"
"				running: true\n"
"				onTriggered: container.destroy();\n"
"			}\n"
"		}\n"
"	}\n"
"\n"
"	function customEmit(x,y) {\n"
"		for (var i=0; i<8; i++) {\n"
"			var obj = emitterComp.createObject(root);\n"
"			obj.x = x\n"
"			obj.y = y\n"
"			obj.targetX = Math.random() * 240 - 120 + obj.x\n"
"			obj.targetY = Math.random() * 240 - 120 + obj.y\n"
"			obj.life = Math.round(Math.random() * 2400) + 200\n"
"			obj.emitRate = Math.round(Math.random() * 32) + 32\n"
"			obj.go();\n"
"		}\n"
"	}\n"
"}\n";

void hack(void *engine_, void *component_)
{
    //static char empty[1] = {0};
    //static char *fakeargv[] = {empty};
    //static int fakeargc = 1;
    //QGuiApplication *app = new QGuiApplication(fakeargc, fakeargv);
    //qApp->setQuitOnLastWindowClosed(false);

    //QQmlEngine *engine = new QQmlEngine(0);
    QQmlEngine *engine = reinterpret_cast<QQmlEngine *>(engine_);
    //QQmlComponent *component = new QQmlComponent(engine);
    QQmlComponent *component = reinterpret_cast<QQmlComponent *>(component_);
    QByteArray ba(particleQml);
    //component->setData(ba, QUrl::fromLocalFile("particle.qml"));
    component->setData(ba, QString("file://./particle.qml"));
    QObject *instance = component->create();
    QQuickView *view = new QQuickView(engine, 0);
    view->setContent(component->url(), component, instance);
    view->setResizeMode(QQuickView::SizeRootObjectToView);
    view->show();
    qApp->exec();
}
}


// vim:ts=4:sw=4:et:ft=cpp

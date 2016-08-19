// #include <private/qmetaobjectbuilder_p.h>

// #include <QtQml/QtQml>
#include <QPainter>
// #include <QQmlEngine>
// #include <QDebug>

#include "painter.h"



void painterSetCompositionMode(QPainter_ *painter, CompositionMode mode) {
  reinterpret_cast<QPainter *>(painter)->setCompositionMode((QPainter::CompositionMode) mode);
}
CompositionMode painterCompositionMode(QPainter_ *painter) {
  return (CompositionMode) reinterpret_cast<QPainter *>(painter)->compositionMode();
}

//
// const QFont &painterFont(QPainter_ *painter) const;
// void painterSetFont(QPainter_ *painter, const QFont &f);
//
// QFontMetrics painterFontMetrics(QPainter_ *painter) const;
// QFontInfo painterFontInfo(QPainter_ *painter) const;
//
// void painterSetPen(QPainter_ *painter, const QColor &color);
// void painterSetPen(QPainter_ *painter, const QPen &pen);
// void painterSetPen(QPainter_ *painter, Qt::PenStyle style);
// const QPen &painterPen(QPainter_ *painter) const;
//
// void painterSetBrush(QPainter_ *painter, const QBrush &brush);
// void painterSetBrush(QPainter_ *painter, Qt::BrushStyle style);
// const QBrush &painterBrush(QPainter_ *painter) const;
//
// // attributes/modes
// void painterSetBackgroundMode(QPainter_ *painter, Qt::BGMode mode);
// Qt::BGMode painterBackgroundMode(QPainter_ *painter) const;
//
// QPoint painterBrushOrigin(QPainter_ *painter) const;
// void painterSetBrushOrigin(QPainter_ *painter, int x, int y);
// void painterSetBrushOrigin(QPainter_ *painter, const QPoint &);
// void painterSetBrushOrigin(QPainter_ *painter, const QPointF &);
//
// void painterSetBackground(QPainter_ *painter, const QBrush &bg);
// const QBrush &painterBackground(QPainter_ *painter) const;
//
// qreal painterOpacity(QPainter_ *painter) const;
// void painterSetOpacity(QPainter_ *painter, qreal opacity);
//
// // Clip functions
// QRegion painterClipRegion(QPainter_ *painter) const;
// QPainterPath painterClipPath(QPainter_ *painter) const;
//
// void painterSetClipRect(QPainter_ *painter, const QRectF &, Qt::ClipOperation op = Qt::ReplaceClip);
// void painterSetClipRect(QPainter_ *painter, const QRect &, Qt::ClipOperation op = Qt::ReplaceClip);
// void painterSetClipRect(QPainter_ *painter, int x, int y, int w, int h, Qt::ClipOperation op = Qt::ReplaceClip);
//
// void painterSetClipRegion(QPainter_ *painter, const QRegion &, Qt::ClipOperation op = Qt::ReplaceClip);
//
// void painterSetClipPath(QPainter_ *painter, const QPainterPath &path, Qt::ClipOperation op = Qt::ReplaceClip);
//
// void painterSetClipping(QPainter_ *painter, bool enable);
// bool painterHasClipping(QPainter_ *painter) const;
//
// QRectF painterClipBoundingRect(QPainter_ *painter) const;
//
void painterSave(QPainter_ *painter){
  reinterpret_cast<QPainter *>(painter)->save();
}
void painterRestore(QPainter_ *painter){
  reinterpret_cast<QPainter *>(painter)->restore();
}
//
// // XForm functions
// void painterSetMatrix(QPainter_ *painter, const QMatrix &matrix, bool combine = false);
// const QMatrix &painterMatrix(QPainter_ *painter) const;
// const QMatrix &painterDeviceMatrix(QPainter_ *painter) const;
// void painterResetMatrix(QPainter_ *painter);
//
// void painterSetTransform(QPainter_ *painter, const QTransform &transform, bool combine = false);
// const QTransform &painterTransform(QPainter_ *painter) const;
// const QTransform &painterDeviceTransform(QPainter_ *painter) const;
// void painterResetTransform(QPainter_ *painter);
//
// void painterSetWorldMatrix(QPainter_ *painter, const QMatrix &matrix, bool combine = false);
// const QMatrix &painterWorldMatrix(QPainter_ *painter) const;
//
// void painterSetWorldTransform(QPainter_ *painter, const QTransform &matrix, bool combine = false);
// const QTransform &painterWorldTransform(QPainter_ *painter) const;
//
// QMatrix painterCombinedMatrix(QPainter_ *painter) const;
// QTransform painterCombinedTransform(QPainter_ *painter) const;
//
// void painterSetMatrixEnabled(QPainter_ *painter, bool enabled);
// bool painterMatrixEnabled(QPainter_ *painter) const;
//
// void painterSetWorldMatrixEnabled(QPainter_ *painter, bool enabled);
// bool painterWorldMatrixEnabled(QPainter_ *painter) const;
//
void painterScale(QPainter_ *painter, qreal sx, qreal sy) {
  reinterpret_cast<QPainter *>(painter)->scale(sx, sy);
}
void painterShear(QPainter_ *painter, qreal sh, qreal sv) {
  reinterpret_cast<QPainter *>(painter)->shear(sh, sv);
}
void painterRotate(QPainter_ *painter, qreal a) {
  reinterpret_cast<QPainter *>(painter)->rotate(a);
}
//
// void painterTranslate(QPainter_ *painter, const QPointF &offset);
// void painterTranslate(QPainter_ *painter, const QPoint &offset);
void painterTranslate(QPainter_ *painter, qreal dx, qreal dy) {
  reinterpret_cast<QPainter *>(painter)->translate(dx, dy);
}
//
// QRect painterWindow(QPainter_ *painter) const;
// void painterSetWindow(QPainter_ *painter, const QRect &window);
void painterSetWindow(QPainter_ *painter, int x, int y, int w, int h) {
  reinterpret_cast<QPainter *>(painter)->setWindow(x, y, w, h);
}
//
// QRect painterViewport(QPainter_ *painter) const;
// void painterSetViewport(QPainter_ *painter, const QRect &viewport);
void painterSetViewport(QPainter_ *painter, int x, int y, int w, int h) {
  reinterpret_cast<QPainter *>(painter)->setViewport(x, y, w, h);
}
//
// void painterSetViewTransformEnabled(QPainter_ *painter, bool enable);
// bool painterViewTransformEnabled(QPainter_ *painter) const;
//
// // drawing functions
// void painterStrokePath(QPainter_ *painter, const QPainterPath &path, const QPen &pen);
// void painterFillPath(QPainter_ *painter, const QPainterPath &path, const QBrush &brush);
// void painterDrawPath(QPainter_ *painter, const QPainterPath &path);
//
// void painterDrawPoint(QPainter_ *painter, const QPointF &pt);
// void painterDrawPoint(QPainter_ *painter, const QPoint &p);
// void painterDrawPoint(QPainter_ *painter, int x, int y);
//
// void painterDrawPoints(QPainter_ *painter, const QPointF *points, int pointCount);
// void painterDrawPoints(QPainter_ *painter, const QPolygonF &points);
// void painterDrawPoints(QPainter_ *painter, const QPoint *points, int pointCount);
// void painterDrawPoints(QPainter_ *painter, const QPolygon &points);
//
// void painterDrawLine(QPainter_ *painter, const QLineF &line);
// void painterDrawLine(QPainter_ *painter, const QLine &line);
// void painterDrawLine(QPainter_ *painter, int x1, int y1, int x2, int y2);
// void painterDrawLine(QPainter_ *painter, const QPoint &p1, const QPoint &p2);
// void painterDrawLine(QPainter_ *painter, const QPointF &p1, const QPointF &p2);
//
// void painterDrawLines(QPainter_ *painter, const QLineF *lines, int lineCount);
// void painterDrawLines(QPainter_ *painter, const QVector<QLineF> &lines);
// void painterDrawLines(QPainter_ *painter, const QPointF *pointPairs, int lineCount);
// void painterDrawLines(QPainter_ *painter, const QVector<QPointF> &pointPairs);
// void painterDrawLines(QPainter_ *painter, const QLine *lines, int lineCount);
// void painterDrawLines(QPainter_ *painter, const QVector<QLine> &lines);
// void painterDrawLines(QPainter_ *painter, const QPoint *pointPairs, int lineCount);
// void painterDrawLines(QPainter_ *painter, const QVector<QPoint> &pointPairs);
//
// void painterDrawRect(QPainter_ *painter, const QRectF &rect);
// void painterDrawRect(QPainter_ *painter, int x1, int y1, int w, int h);
// void painterDrawRect(QPainter_ *painter, const QRect &rect);
//
// void painterDrawRects(QPainter_ *painter, const QRectF *rects, int rectCount);
// void painterDrawRects(QPainter_ *painter, const QVector<QRectF> &rectangles);
// void painterDrawRects(QPainter_ *painter, const QRect *rects, int rectCount);
// void painterDrawRects(QPainter_ *painter, const QVector<QRect> &rectangles);
//
// void painterDrawEllipse(QPainter_ *painter, const QRectF &r);
// void painterDrawEllipse(QPainter_ *painter, const QRect &r);
// void painterDrawEllipse(QPainter_ *painter, int x, int y, int w, int h);
//
// void painterDrawEllipse(QPainter_ *painter, const QPointF &center, qreal rx, qreal ry);
// void painterDrawEllipse(QPainter_ *painter, const QPoint &center, int rx, int ry);
//
// void painterDrawPolyline(QPainter_ *painter, const QPointF *points, int pointCount);
// void painterDrawPolyline(QPainter_ *painter, const QPolygonF &polyline);
// void painterDrawPolyline(QPainter_ *painter, const QPoint *points, int pointCount);
// void painterDrawPolyline(QPainter_ *painter, const QPolygon &polygon);
//
// void painterDrawPolygon(QPainter_ *painter, const QPointF *points, int pointCount, Qt::FillRule fillRule = Qt::OddEvenFill);
// void painterDrawPolygon(QPainter_ *painter, const QPolygonF &polygon, Qt::FillRule fillRule = Qt::OddEvenFill);
// void painterDrawPolygon(QPainter_ *painter, const QPoint *points, int pointCount, Qt::FillRule fillRule = Qt::OddEvenFill);
// void painterDrawPolygon(QPainter_ *painter, const QPolygon &polygon, Qt::FillRule fillRule = Qt::OddEvenFill);
//
// void painterDrawConvexPolygon(QPainter_ *painter, const QPointF *points, int pointCount);
// void painterDrawConvexPolygon(QPainter_ *painter, const QPolygonF &polygon);
// void painterDrawConvexPolygon(QPainter_ *painter, const QPoint *points, int pointCount);
// void painterDrawConvexPolygon(QPainter_ *painter, const QPolygon &polygon);
//
// void painterDrawArc(QPainter_ *painter, const QRectF &rect, int a, int alen);
// void painterDrawArc(QPainter_ *painter, const QRect &, int a, int alen);
// void painterDrawArc(QPainter_ *painter, int x, int y, int w, int h, int a, int alen);
//
// void painterDrawPie(QPainter_ *painter, const QRectF &rect, int a, int alen);
// void painterDrawPie(QPainter_ *painter, int x, int y, int w, int h, int a, int alen);
// void painterDrawPie(QPainter_ *painter, const QRect &, int a, int alen);
//
// void painterDrawChord(QPainter_ *painter, const QRectF &rect, int a, int alen);
// void painterDrawChord(QPainter_ *painter, int x, int y, int w, int h, int a, int alen);
// void painterDrawChord(QPainter_ *painter, const QRect &, int a, int alen);
//
// void painterDrawRoundedRect(QPainter_ *painter, const QRectF &rect, qreal xRadius, qreal yRadius,
//                      Qt::SizeMode mode = Qt::AbsoluteSize);
// void painterDrawRoundedRect(QPainter_ *painter, int x, int y, int w, int h, qreal xRadius, qreal yRadius,
//                             Qt::SizeMode mode = Qt::AbsoluteSize);
// void painterDrawRoundedRect(QPainter_ *painter, const QRect &rect, qreal xRadius, qreal yRadius,
//                             Qt::SizeMode mode = Qt::AbsoluteSize);
//
// void painterDrawRoundRect(QPainter_ *painter, const QRectF &r, int xround = 25, int yround = 25);
// void painterDrawRoundRect(QPainter_ *painter, int x, int y, int w, int h, int = 25, int = 25);
// void painterDrawRoundRect(QPainter_ *painter, const QRect &r, int xround = 25, int yround = 25);
//
// void painterDrawTiledPixmap(QPainter_ *painter, const QRectF &rect, const QPixmap &pm, const QPointF &offset = QPointF());
// void painterDrawTiledPixmap(QPainter_ *painter, int x, int y, int w, int h, const QPixmap &, int sx=0, int sy=0);
// void painterDrawTiledPixmap(QPainter_ *painter, const QRect &, const QPixmap &, const QPoint & = QPoint());
// #ifndef QT_NO_PICTURE
// void painterDrawPicture(QPainter_ *painter, const QPointF &p, const QPicture &picture);
// void painterDrawPicture(QPainter_ *painter, int x, int y, const QPicture &picture);
// void painterDrawPicture(QPainter_ *painter, const QPoint &p, const QPicture &picture);
// #endif
//
// void painterDrawPixmap(QPainter_ *painter, const QRectF &targetRect, const QPixmap &pixmap, const QRectF &sourceRect);
// void painterDrawPixmap(QPainter_ *painter, const QRect &targetRect, const QPixmap &pixmap, const QRect &sourceRect);
// void painterDrawPixmap(QPainter_ *painter, int x, int y, int w, int h, const QPixmap &pm,
//                        int sx, int sy, int sw, int sh);
// void painterDrawPixmap(QPainter_ *painter, int x, int y, const QPixmap &pm,
//                        int sx, int sy, int sw, int sh);
// void painterDrawPixmap(QPainter_ *painter, const QPointF &p, const QPixmap &pm, const QRectF &sr);
// void painterDrawPixmap(QPainter_ *painter, const QPoint &p, const QPixmap &pm, const QRect &sr);
// void painterDrawPixmap(QPainter_ *painter, const QPointF &p, const QPixmap &pm);
// void painterDrawPixmap(QPainter_ *painter, const QPoint &p, const QPixmap &pm);
// void painterDrawPixmap(QPainter_ *painter, int x, int y, const QPixmap &pm);
// void painterDrawPixmap(QPainter_ *painter, const QRect &r, const QPixmap &pm);
// void painterDrawPixmap(QPainter_ *painter, int x, int y, int w, int h, const QPixmap &pm);
//
// void painterDrawPixmapFragments(QPainter_ *painter, const PixmapFragment *fragments, int fragmentCount,
//                          const QPixmap &pixmap, PixmapFragmentHints hints = 0);
//
// void painterDrawImage(QPainter_ *painter, const QRectF &targetRect, const QImage &image, const QRectF &sourceRect,
//                Qt::ImageConversionFlags flags = Qt::AutoColor);
// void painterDrawImage(QPainter_ *painter, const QRect &targetRect, const QImage &image, const QRect &sourceRect,
//                       Qt::ImageConversionFlags flags = Qt::AutoColor);
// void painterDrawImage(QPainter_ *painter, const QPointF &p, const QImage &image, const QRectF &sr,
//                       Qt::ImageConversionFlags flags = Qt::AutoColor);
// void painterDrawImage(QPainter_ *painter, const QPoint &p, const QImage &image, const QRect &sr,
//                       Qt::ImageConversionFlags flags = Qt::AutoColor);
// void painterDrawImage(QPainter_ *painter, const QRectF &r, const QImage &image);
// void painterDrawImage(QPainter_ *painter, const QRect &r, const QImage &image);
// void painterDrawImage(QPainter_ *painter, const QPointF &p, const QImage &image);
// void painterDrawImage(QPainter_ *painter, const QPoint &p, const QImage &image);
// void painterDrawImage(QPainter_ *painter, int x, int y, const QImage &image, int sx = 0, int sy = 0,
//                       int sw = -1, int sh = -1, Qt::ImageConversionFlags flags = Qt::AutoColor);
//
// void painterSetLayoutDirection(QPainter_ *painter, Qt::LayoutDirection direction);
// Qt::LayoutDirection painterLayoutDirection(QPainter_ *painter) const;
//
// #if !defined(QT_NO_RAWFONT)
// void painterDrawGlyphRun(QPainter_ *painter, const QPointF &position, const QGlyphRun &glyphRun);
// #endif
//
// void painterDrawStaticText(QPainter_ *painter, const QPointF &topLeftPosition, const QStaticText &staticText);
// void painterDrawStaticText(QPainter_ *painter, const QPoint &topLeftPosition, const QStaticText &staticText);
// void painterDrawStaticText(QPainter_ *painter, int left, int top, const QStaticText &staticText);
//
// void painterDrawText(QPainter_ *painter, const QPointF &p, const QString &s);
// void painterDrawText(QPainter_ *painter, const QPoint &p, const QString &s);
// void painterDrawText(QPainter_ *painter, int x, int y, const QString &s);
//
// void painterDrawText(QPainter_ *painter, const QPointF &p, const QString &str, int tf, int justificationPadding);
//
// void painterDrawText(QPainter_ *painter, const QRectF &r, int flags, const QString &text, QRectF *br=0);
// void painterDrawText(QPainter_ *painter, const QRect &r, int flags, const QString &text, QRect *br=0);
// void painterDrawText(QPainter_ *painter, int x, int y, int w, int h, int flags, const QString &text, QRect *br=0);
//
// void painterDrawText(QPainter_ *painter, const QRectF &r, const QString &text, const QTextOption &o = QTextOption());
//
// QRectF painterBoundingRect(QPainter_ *painter, const QRectF &rect, int flags, const QString &text);
// QRect painterBoundingRect(QPainter_ *painter, const QRect &rect, int flags, const QString &text);
// QRect painterBoundingRect(QPainter_ *painter, int x, int y, int w, int h, int flags, const QString &text);
//
// QRectF painterBoundingRect(QPainter_ *painter, const QRectF &rect, const QString &text, const QTextOption &o = QTextOption());
//
// void painterDrawTextItem(QPainter_ *painter, const QPointF &p, const QTextItem &ti);
// void painterDrawTextItem(QPainter_ *painter, int x, int y, const QTextItem &ti);
// void painterDrawTextItem(QPainter_ *painter, const QPoint &p, const QTextItem &ti);
//
// void painterFillRect(QPainter_ *painter, const QRectF &, const QBrush &);
// void painterFillRect(QPainter_ *painter, int x, int y, int w, int h, const QBrush &);
// void painterFillRect(QPainter_ *painter, const QRect &, const QBrush &);
//
// void painterFillRect(QPainter_ *painter, const QRectF &, const QColor &color);
// void painterFillRect(QPainter_ *painter, int x, int y, int w, int h, const QColor &color);
// void painterFillRect(QPainter_ *painter, const QRect &, const QColor &color);
//
// void painterFillRect(QPainter_ *painter, int x, int y, int w, int h, Qt::GlobalColor c);
// void painterFillRect(QPainter_ *painter, const QRect &r, Qt::GlobalColor c);
// void painterFillRect(QPainter_ *painter, const QRectF &r, Qt::GlobalColor c);
//
// void painterFillRect(QPainter_ *painter, int x, int y, int w, int h, Qt::BrushStyle style);
// void painterFillRect(QPainter_ *painter, const QRect &r, Qt::BrushStyle style);
// void painterFillRect(QPainter_ *painter, const QRectF &r, Qt::BrushStyle style);
//
// void painterEraseRect(QPainter_ *painter, const QRectF &);
// void painterEraseRect(QPainter_ *painter, int x, int y, int w, int h);
// void painterEraseRect(QPainter_ *painter, const QRect &);
//
// void painterSetRenderHint(QPainter_ *painter, RenderHint hint, bool on = true);
// void painterSetRenderHints(QPainter_ *painter, RenderHints hints, bool on = true);
// RenderHints painterRenderHints(QPainter_ *painter) const;
// bool painterTestRenderHint(QPainter_ *painter, RenderHint hint) const { return renderHints() & hint; }
//
// QPaintEngine *painterPaintEngine(QPainter_ *painter) const;

// static void setRedirected(const QPaintDevice *device, QPaintDevice *replacement,
//                           const QPoint& offset = QPoint());
// static QPaintDevice *redirected(const QPaintDevice *device, QPoint *offset = 0);
// static void restoreRedirected(const QPaintDevice *device);

void painterBeginNativePainting(QPainter_ *painter) {
  reinterpret_cast<QPainter *>(painter)->beginNativePainting();
}
void painterEndNativePainting(QPainter_ *painter) {
  reinterpret_cast<QPainter *>(painter)->endNativePainting();
}

// vim:ts=4:sw=4:et:ft=cpp

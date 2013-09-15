# Maintainer: ushi <ushi@honkgong.info>
pkgname='resize-git'
pkgver=3.30cd9d4
pkgrel=1
pkgdesc='image resizing server'
arch=('x86_64')
url='https://resize.honkgong.info'
license=('MIT')
conflicts=('resize')
provides=('resize')
makedepends=('go')
source=('resize::git+https://github.com/ushis/resize-demo#branch=master')
sha256sums=('SKIP')

pkgver() {
  cd resize
  echo "$(git rev-list --count master).$(git rev-parse --short master)"
}

build() {
  cd resize
  make
}

package() {
  cd resize
  install -Dm0755 resized                 "${pkgdir}/usr/bin/resized"
  install -dm0755                         "${pkgdir}/usr/share/resize"
  install -Dm0644 root/*                  "${pkgdir}/usr/share/resize"
  install -Dm0644 systemd/resized.conf    "${pkgdir}/usr/lib/tmpfiles.d/resized.conf"
  install -Dm0644 systemd/resized.service "${pkgdir}/usr/lib/systemd/system/resized.service"
}

# vim:set ts=2 sw=2 et:

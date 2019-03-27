#!/usr/bin/env bash

pkg_version="${PKG_VERSION}"
tag_tag="${TAG_NAME}"

if [ ! "${pkg_version}" ] && [ ! "${tag_tag}" ]; then
    echo "No version or tag specified."
    exit 1
fi

if [ "${pkg_version}" != "${tag_tag}" ]; then
    echo "Versions do not match: pkg@${pkg_version} tag@${tag_tag}"
    exit 1
fi

echo "Versions match: pkg@${pkg_version} tag@${tag_tag}"

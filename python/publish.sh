uv version --bump patch
rm -rf ./dist
uv version
uv build
uv publish ./dist/* --index minieye-pypi

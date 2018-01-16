#!/bin/bash

if [ "$1" == "onlp" ]; then
	cd pluginManager/onlp/libs
	ln -sf libonlp.so.1 libonlp.so
	cd -
elif [ "$1" == "openBMC" ]; then
	echo "Build Target is OpenBMC"
elif [ "$1" == "openBMCVoyager" ]; then
	echo "Build Target is OpenBMCVoyager"
else
	echo "Build Target is Dummy"
fi

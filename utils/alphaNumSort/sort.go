//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package alphaNumSort

import (
	"math"
	"sort"
)

const (
	MAX_CHAR_ASCII_VAL = 123 //ASCII value of 'z' + 1
)

/* Simple sort routine that sorts alpha numeric strings
   containing the following runes (0-9, a-z, A-Z, _).
*/
func Sort(strList []string) []string {
	if (strList == nil) || (len(strList) == 1) {
		return strList
	}

	var outList []string = make([]string, len(strList))
	var strMap map[float64]string = make(map[float64]string)
	for _, str := range strList {
		wt := computeWeight(str)
		strMap[wt] = str
	}
	keySlice := make([]float64, len(strMap))
	idx := 0
	for key, _ := range strMap {
		keySlice[idx] = key
		idx++
	}
	sort.Float64s(keySlice)
	for idx, key := range keySlice {
		outList[idx] = strMap[key]
	}
	return outList
}

func Compare(s1, s2 string) int {
	wt1 := computeWeight(s1)
	wt2 := computeWeight(s2)
	if wt1 < wt2 {
		return -1
	} else if wt1 == wt2 {
		return 0
	} else {
		return 1
	}
}

/* Computes weight of give string. Max char ascii val ('z') = 122 */
func computeWeight(str string) float64 {
	var wt float64
	for idx, val := range str {
		wt += math.Pow(float64(MAX_CHAR_ASCII_VAL), float64(idx)) * float64(val)
	}
	return wt
}

package api

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Helper Functions

// Sets a key in a map to n, or increments its value by n if key exists.
func setOrIncrement(m map[int]int, key int, n int) map[int]int {
	_, ok := m[key]
	if !ok {
		m[key] = n
	} else {
		m[key] += n
	}
	return m
}

// Generates an optimized shipment summary for given order quantity and pack sizes
func generateShipment(q int, sizes []int) (map[int]int, int) {
	// Set the iterator `i` to the last index of `sizes`
	i := len(sizes) - 1
	res := map[int]int{}
	// While order quantity is greater than zero, pick packages appropriately
	for q > 0 {
		// Current pack size is `s`
		s := sizes[i]
		// If picking `s` would conclude the shipment,
		if q-s < 0 {
			// If `s` is the smallest pack available, add it to the order and
			// close the shipment.
			if i == 0 {
				res = setOrIncrement(res, s, 1)
				q -= s
				break
			} else {
				// In order for deciding if we should pick `s`, we check what
				// would be the excess items if we didn't pick `s`.
				_, _q := generateShipment(q, sizes[:i])

				// `q-s` is the excess amount given `s` is picked.
				// _q is the excess amount given `s` is NOT picked.
				// If excess amount when we pick `s` is smaller than otherwise,
				// we pick `s` and conclude the order.
				if q-s >= _q {
					res = setOrIncrement(res, s, 1)
					q -= s
					break
				}

				// Otherwise, we just skip `s`.
				i--
				continue
			}
		}

		// If picking `s` would not conclude the shipment,
		// we deduct the most possible amount of items using `s`
		// from the total order count.
		// This way, we ensure that minimum number of packs are
		// being used.
		n := q / s
		res = setOrIncrement(res, s, n)
		q -= n * s
	}
	return res, q
}

// Helper function to read array `sizes` from context and convert its elements to integer
func getSizesFromContext(c *gin.Context) ([]int, error) {
	__sizes := c.Query("sizes")
	if __sizes == "" {
		if os.Getenv("DEFAULT_PACK_SIZES") != "" {
			__sizes = os.Getenv("DEFAULT_PACK_SIZES")
		} else {
			return []int{250, 500, 1000, 2000, 5000}, nil
		}
	}
	_sizes := strings.Split(__sizes, ",")
	var sizes []int
	for _, _s := range _sizes {
		s, err := strconv.Atoi(strings.TrimSpace(_s))
		if err != nil {
			return []int{}, err
		}
		sizes = append(sizes, s)
	}
	return sizes, nil
}

// Wrapper function to get query params `sizes` and `quantity` from context and create an
// order summary out of them.
func generateShipmentFromContext(c *gin.Context) (map[int]int, int, []int, error) {
	var res map[int]int
	v, err := strconv.Atoi(c.Query("quantity"))

	if err != nil || v < 0 || v > 2147483647 {
		return res, v, []int{}, errors.New("Cannot parse quantity. Please check your input.")
	}

	sizes, err := getSizesFromContext(c)

	if err != nil {
		return res, v, []int{}, errors.New("Cannot parse sizes. Please check your input.")
	}

	slices.Sort(sizes)

	res, _ = generateShipment(v, sizes)
	return res, v, sizes, nil
}

// Views

// Returns a static index HTML page
func IndexHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/index.tmpl", gin.H{})
}

// Returns a HTML page including the requested order's summary
func GetShipmentHTML(c *gin.Context) {
	res, v, sizes, err := generateShipmentFromContext(c)

	c.HTML(http.StatusOK, "pages/order.tmpl", gin.H{
		"shipmentSummary": res,
		"err":             err,
		"quantity":        v,
		"sizes":           sizes,
	})
}

// Returns a JSON including the requested order's summary
func GetShipment(c *gin.Context) {
	res, _, _, err := generateShipmentFromContext(c)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

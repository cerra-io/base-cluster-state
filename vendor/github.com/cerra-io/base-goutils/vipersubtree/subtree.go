/*
 * This is a poor-man replacement of the viper.Sub() subtree feature of
 * viper. The latter does not properly handle overrides (from environment vars
 * or from command line) in that it ignores any such overrides one a subtree
 * is requested with the Sub(prefix) method. We consider it a bug in viper.
 *
 * This workaround provides a subtree by calling st := Subtree(viper, prefix). All
 * viper.Get*() methods work on the subtree and properly return the overridden
 * values.
 */

package vipersubtree

import (
	"time"

	"github.com/spf13/viper"
)

// ViperSubtree defines a self-sufficient subtree of config entries
type ViperSubtree struct {
	*viper.Viper
	prefix string
}

// Subtree returns a subtree
func Subtree(conf *viper.Viper, prefix string) *ViperSubtree {
	vst := &ViperSubtree{conf, prefix}
	return vst
}

// Get is same as Viper.Get
func (st *ViperSubtree) Get(key string) interface{} {
	return st.Viper.Get(st.prefix + "." + key)
}

// GetBool is same as Viper.GetBool
func (st *ViperSubtree) GetBool(key string) bool {
	return st.Viper.GetBool(st.prefix + "." + key)
}

// GetDuration is same as Viper.GetDuration
func (st *ViperSubtree) GetDuration(key string) time.Duration {
	return st.Viper.GetDuration(st.prefix + "." + key)
}

// GetFloat64 is same as Viper.GetFloat64
func (st *ViperSubtree) GetFloat64(key string) float64 {
	return st.Viper.GetFloat64(st.prefix + "." + key)
}

// GetInt is same as Viper.GetInt
func (st *ViperSubtree) GetInt(key string) int {
	return st.Viper.GetInt(st.prefix + "." + key)
}

// GetInt64 is same as Viper.GetInt64
func (st *ViperSubtree) GetInt64(key string) int64 {
	return st.Viper.GetInt64(st.prefix + "." + key)
}

// GetSizeInBytes is same as Viper.GetSizeInBytes
func (st *ViperSubtree) GetSizeInBytes(key string) uint {
	return st.Viper.GetSizeInBytes(st.prefix + "." + key)
}

// GetString is same as Viper.GetString
func (st *ViperSubtree) GetString(key string) string {
	return st.Viper.GetString(st.prefix + "." + key)
}

// GetStringMap is same as Viper.GetStringMap
func (st *ViperSubtree) GetStringMap(key string) map[string]interface{} {
	return st.Viper.GetStringMap(st.prefix + "." + key)
}

// GetStringMapString is same as Viper.GetStringMapString
func (st *ViperSubtree) GetStringMapString(key string) map[string]string {
	return st.Viper.GetStringMapString(st.prefix + "." + key)
}

// GetStringMapStringSlice is same as Viper.GetStringMapStringSlice
func (st *ViperSubtree) GetStringMapStringSlice(key string) map[string][]string {
	return st.Viper.GetStringMapStringSlice(st.prefix + "." + key)
}

// GetStringSlice is same as Viper.GetStringSlice
func (st *ViperSubtree) GetStringSlice(key string) []string {
	return st.Viper.GetStringSlice(st.prefix + "." + key)
}

// GetTime is same as Viper.GetTime
func (st *ViperSubtree) GetTime(key string) time.Time {
	return st.Viper.GetTime(st.prefix + "." + key)
}

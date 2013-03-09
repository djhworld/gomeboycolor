/**
 * Created with IntelliJ IDEA.
 * User: danielharper
 * Date: 01/03/2013
 * Time: 21:39
 * To change this template use File | Settings | File Templates.
 */
package components

type IRQHandler interface {
	RequestInterrupt(interrupt byte)
}


/**
 * Base element of HTD Library Editor
 */
export default class HtdElement {
  constructor(parent, target) {
    /**
     * @type {HtdElement} element
     * @public
     */
    this.parent = parent;

    /**
     * @type {HtdElement[]} children
     * @public
     */
    this.children = [];

    /**
     * @type {HTMLElement} element
     * @public
     */
    this.target = target;

    /**
     * @type {HTMLElement} element
     * @public
     */
    this.element;
  }

  /**
   * @param {HtdElement} child
   */
  addChild(child) {
    this.children = [...this.children, child];
  }

  /**
   * @param {number} index
   */
  getChild(index) {
    return this.children[index];
  }
}

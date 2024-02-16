import HtdElement from "./HtdElement";

export default class Container {
  /**
   * @param {HtdElement} parent
   * @param {HTMLElement} target
   */
  constructor(parent, target) {
    this.element = document.createElement("div");
    this.target = target;
    this.parent = parent;
  }

  getElement() { return this.container; }
  getTarget() { return this.target; }
  getParent() { return this.parent; }

  /**
   * @param {HtdElement} element
   */
  addChild(element) {

  }

  /**
   * @param {HtdElement} child 
   */
  removeChild(child) {


  }
}

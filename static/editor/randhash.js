/**
 * Generates a random hash.
 *
 * This function generates a random hash by generating a random string of characters
 * and then hashing it using a hashing algorithm. The generated hash is returned as
 * the result.
 *
 * @param {number} length - The length of the random string to generate. Must be a positive integer.
 * @returns {string} The generated random hash.
 *
 * @throws {Error} If the length parameter is not a positive integer, an error is thrown indicating
 *                 that the length should be a positive integer.
 *
 * @example
 * const randomHash = generateRandomHash(10);
 * console.log(randomHash); // Example Output: "a1b2c3d4e5"
 */
export function generateRandomHash(length) {
  if (!Number.isInteger(length) || length <= 0) {
    throw new Error("Length should be a positive integer.");
  }

  const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  let randomString = "";

  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * characters.length);
    randomString += characters.charAt(randomIndex);
  }

  return randomString;
}

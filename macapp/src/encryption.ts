import { createCipheriv, createDecipheriv, randomBytes } from 'crypto'
import { exec } from 'child_process'
import { promisify } from 'util'

const execAsync = promisify(exec)

const ALGORITHM = 'aes-256-gcm'
const KEY_SIZE = 32 // 256 bits
const IV_SIZE = 12 // GCM standard
const AUTH_TAG_SIZE = 16

// Get or create encryption key from macOS Keychain
async function getEncryptionKey(): Promise<Buffer> {
  try {
    // Try to get existing key
    const { stdout } = await execAsync('security find-generic-password -s "secllama-encryption-key" -w 2>/dev/null')
    if (stdout && stdout.trim()) {
      const keyHex = stdout.trim()
      // Validate it's a proper hex string
      if (/^[0-9a-f]+$/i.test(keyHex)) {
        return Buffer.from(keyHex, 'hex')
      }
    }
  } catch (err) {
    // Key doesn't exist, create it
    console.log('Creating new encryption key...')
  }

  // Generate new key (32 bytes for AES-256)
  const key = randomBytes(KEY_SIZE)
  const keyHex = key.toString('hex')

  try {
    // Delete old key if exists, then add new one
    await execAsync(`security delete-generic-password -s "secllama-encryption-key" 2>/dev/null || true`)
    await execAsync(
      `security add-generic-password -s "secllama-encryption-key" -a "secllama" -w "${keyHex}" -T /usr/bin/security`
    )
    console.log('✓ Encryption key created and stored in macOS Keychain')
  } catch (err) {
    console.error('Failed to store key in keychain:', err)
    // Continue anyway - key is in memory
  }

  return key
}

// Ensure encryption key exists (call this on app startup)
export async function initializeEncryption(): Promise<void> {
  try {
    await getEncryptionKey()
    console.log('✓ Encryption initialized')
  } catch (err) {
    console.error('Failed to initialize encryption:', err)
  }
}

// Encrypt data using AES-256-GCM (same as Go implementation)
export async function encryptData(plaintext: string): Promise<string> {
  try {
    const key = await getEncryptionKey()
    const iv = randomBytes(IV_SIZE)
    
    const cipher = createCipheriv(ALGORITHM, key, iv)
    
    let encrypted = cipher.update(plaintext, 'utf8')
    encrypted = Buffer.concat([encrypted, cipher.final()])
    
    const authTag = cipher.getAuthTag()
    
    // Combine: iv + encrypted + authTag
    const combined = Buffer.concat([iv, encrypted, authTag])
    
    return combined.toString('base64')
  } catch (err) {
    console.error('Encryption failed:', err)
    throw err
  }
}

// Decrypt data using AES-256-GCM
export async function decryptData(ciphertext: string): Promise<string> {
  try {
    const key = await getEncryptionKey()
    const combined = Buffer.from(ciphertext, 'base64')
    
    // Extract: iv + encrypted + authTag
    const iv = combined.slice(0, IV_SIZE)
    const authTag = combined.slice(-AUTH_TAG_SIZE)
    const encrypted = combined.slice(IV_SIZE, -AUTH_TAG_SIZE)
    
    const decipher = createDecipheriv(ALGORITHM, key, iv)
    decipher.setAuthTag(authTag)
    
    let decrypted = decipher.update(encrypted)
    decrypted = Buffer.concat([decrypted, decipher.final()])
    
    return decrypted.toString('utf8')
  } catch (err) {
    console.error('Decryption failed:', err)
    throw err
  }
}

// Encrypt conversations before storing
export async function encryptConversations(conversations: any[]): Promise<string> {
  const json = JSON.stringify(conversations)
  return await encryptData(json)
}

// Decrypt conversations after loading
export async function decryptConversations(encrypted: string): Promise<any[]> {
  try {
    const json = await decryptData(encrypted)
    return JSON.parse(json)
  } catch (err) {
    console.error('Failed to decrypt conversations:', err)
    return []
  }
}


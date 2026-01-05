import { ElMessage } from 'element-plus'

/**
 * 复制文本到剪贴板
 * @param text 要复制的文本
 * @returns Promise<boolean> 成功返回true，失败返回false
 */
export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
    return true
  } catch (err) {
    console.error('复制失败:', err)
    ElMessage.error('复制失败')
    return false
  }
}

/**
 * 兼容旧浏览器的复制方法
 * @param text 要复制的文本
 */
export const fallbackCopyToClipboard = (text: string): boolean => {
  const textArea = document.createElement('textarea')
  textArea.value = text
  textArea.style.position = 'fixed'
  textArea.style.top = '0'
  textArea.style.left = '0'
  textArea.style.width = '2em'
  textArea.style.height = '2em'
  textArea.style.padding = '0'
  textArea.style.border = 'none'
  textArea.style.outline = 'none'
  textArea.style.boxShadow = 'none'
  textArea.style.background = 'transparent'
  document.body.appendChild(textArea)
  textArea.focus()
  textArea.select()

  try {
    const successful = document.execCommand('copy')
    document.body.removeChild(textArea)
    if (successful) {
      ElMessage.success('已复制到剪贴板')
    } else {
      ElMessage.error('复制失败')
    }
    return successful
  } catch (err) {
    document.body.removeChild(textArea)
    ElMessage.error('复制失败')
    return false
  }
}

/**
 * 智能复制文本（自动选择最佳方法）
 * @param text 要复制的文本
 * @returns Promise<boolean> 成功返回true，失败返回false
 */
export const smartCopy = async (text: string): Promise<boolean> => {
  if (navigator.clipboard) {
    return await copyToClipboard(text)
  }
  return fallbackCopyToClipboard(text)
}

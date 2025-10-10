(function() {
  const API_BASE = window.COMMENT_API_BASE || 'http://localhost:8080';
  const ARTICLE_URL = window.location.href.split('#')[0].split('?')[0];

  const styles = `
    .icomment-container { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; width: 100%; }
    .icomment-form { margin-bottom: 30px; padding: 20px; border: 1px solid #e0e0e0; border-radius: 4px; }
    .icomment-input { width: 100%; padding: 8px; margin-bottom: 10px; border: 1px solid #ddd; border-radius: 3px; box-sizing: border-box; }
    .icomment-textarea { width: 100%; padding: 8px; margin-bottom: 10px; border: 1px solid #ddd; border-radius: 3px; min-height: 100px; box-sizing: border-box; }
    .icomment-button { padding: 10px 20px; background: #333; color: white; border: none; border-radius: 3px; cursor: pointer; }
    .icomment-button:hover { background: #555; }
    .icomment-list { list-style: none; padding: 0; }
    .icomment-item { margin-bottom: 20px; padding: 15px; border: 1px solid #e0e0e0; border-radius: 4px; }
    .icomment-reply { margin-left: 40px; margin-top: 10px; padding: 10px; background: #f9f9f9; border-left: 3px solid #ddd; }
    .icomment-meta { font-size: 0.9em; color: #666; margin-bottom: 8px; }
    .icomment-content { margin: 10px 0; line-height: 1.6; }
    .icomment-reply-btn { font-size: 0.85em; color: #666; cursor: pointer; text-decoration: underline; }
    .icomment-cancel-btn { margin-left: 10px; font-size: 0.85em; color: #999; cursor: pointer; text-decoration: underline; }
  `;

  const styleSheet = document.createElement('style');
  styleSheet.textContent = styles;
  document.head.appendChild(styleSheet);

  // Validation constants
  const MAX_NICKNAME_LENGTH = 50;
  const MAX_EMAIL_LENGTH = 100;
  const MAX_CONTENT_LENGTH = 2000;

  function validateEmail(email) {
    if (!email) return true; // Email is optional
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  }

  function validateInput(nickname, email, content) {
    if (!nickname || nickname.length > MAX_NICKNAME_LENGTH) {
      alert(`Nickname is required and must be less than ${MAX_NICKNAME_LENGTH} characters.`);
      return false;
    }
    
    if (email && !validateEmail(email)) {
      alert('Please enter a valid email address.');
      return false;
    }
    
    if (email && email.length > MAX_EMAIL_LENGTH) {
      alert(`Email must be less than ${MAX_EMAIL_LENGTH} characters.`);
      return false;
    }
    
    if (!content || content.length > MAX_CONTENT_LENGTH) {
      alert(`Comment is required and must be less than ${MAX_CONTENT_LENGTH} characters.`);
      return false;
    }
    
    return true;
  }

  function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  function formatDate(dateStr) {
    if (!dateStr) return '';
    // SQLite returns: 2025-10-10 12:34:56
    // Convert to ISO format for parsing
    const isoDate = dateStr.replace(' ', 'T') + 'Z';
    const date = new Date(isoDate);
    if (isNaN(date.getTime())) return dateStr;
    return date.toLocaleString();
  }

  function renderComments(comments) {
    const list = document.getElementById('icomment-list');
    list.innerHTML = '';

    const topLevel = comments.filter(c => !c.parent_id);
    const replies = comments.filter(c => c.parent_id);

    topLevel.forEach(comment => {
      const item = document.createElement('div');
      item.className = 'icomment-item';
      item.innerHTML = `
        <div class="icomment-meta">
          <strong>${escapeHtml(comment.nickname)}</strong> | ${formatDate(comment.created_at)}
        </div>
        <div class="icomment-content">${escapeHtml(comment.content)}</div>
        <span class="icomment-reply-btn" data-id="${comment.id}">Reply</span>
        <div id="reply-form-${comment.id}"></div>
      `;

      const commentReplies = replies.filter(r => r.parent_id === comment.id);
      commentReplies.forEach(reply => {
        const replyDiv = document.createElement('div');
        replyDiv.className = 'icomment-reply';
        replyDiv.innerHTML = `
          <div class="icomment-meta">
            <strong>${escapeHtml(reply.nickname)}</strong> | ${formatDate(reply.created_at)}
          </div>
          <div class="icomment-content">${escapeHtml(reply.content)}</div>
        `;
        item.appendChild(replyDiv);
      });

      list.appendChild(item);
    });

    document.querySelectorAll('.icomment-reply-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const parentId = e.target.dataset.id;
        showReplyForm(parentId);
      });
    });
  }

  function showReplyForm(parentId) {
    const container = document.getElementById(`reply-form-${parentId}`);
    container.innerHTML = `
      <div style="margin-top: 10px; padding: 10px; background: #f5f5f5; border-radius: 3px;">
        <input type="text" id="reply-nickname-${parentId}" class="icomment-input" placeholder="Nickname" maxlength="${MAX_NICKNAME_LENGTH}" required>
        <input type="email" id="reply-email-${parentId}" class="icomment-input" placeholder="Email (optional)" maxlength="${MAX_EMAIL_LENGTH}">
        <textarea id="reply-content-${parentId}" class="icomment-textarea" placeholder="Your reply..." maxlength="${MAX_CONTENT_LENGTH}" required></textarea>
        <button class="icomment-button" onclick="submitReply(${parentId})">Submit Reply</button>
        <span class="icomment-cancel-btn" onclick="cancelReply(${parentId})">Cancel</span>
      </div>
    `;
  }

  window.submitReply = function(parentId) {
    const nickname = document.getElementById(`reply-nickname-${parentId}`).value.trim();
    const email = document.getElementById(`reply-email-${parentId}`).value.trim();
    const content = document.getElementById(`reply-content-${parentId}`).value.trim();

    if (!validateInput(nickname, email, content)) {
      return;
    }

    fetch(`${API_BASE}/api/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        article_url: ARTICLE_URL,
        parent_id: parentId,
        nickname,
        email,
        content
      })
    })
    .then(res => {
      if (!res.ok) {
        return res.json().then(err => Promise.reject(err));
      }
      return res.json();
    })
    .then(data => {
      alert('Reply submitted successfully! It will appear after approval.');
      loadComments();
      cancelReply(parentId);
    })
    .catch(err => {
      const message = err.error || 'Failed to submit reply';
      alert(message + '\n\nPlease try refreshing the page and submitting again.');
    });
  };

  window.cancelReply = function(parentId) {
    document.getElementById(`reply-form-${parentId}`).innerHTML = '';
  };

  function loadComments() {
    fetch(`${API_BASE}/api/comments?article_url=${encodeURIComponent(ARTICLE_URL)}`)
      .then(res => res.json())
      .then(comments => renderComments(comments || []))
      .catch(err => console.error('Error loading comments:', err));
  }

  function init() {
    const container = document.getElementById('icomment');
    if (!container) return;

    container.innerHTML = `
      <div class="icomment-container">
        <div class="icomment-form">
          <h3>Leave a Comment</h3>
          <input type="text" id="icomment-nickname" class="icomment-input" placeholder="Nickname" maxlength="${MAX_NICKNAME_LENGTH}" required>
          <input type="email" id="icomment-email" class="icomment-input" placeholder="Email (optional)" maxlength="${MAX_EMAIL_LENGTH}">
          <textarea id="icomment-content" class="icomment-textarea" placeholder="Your comment..." maxlength="${MAX_CONTENT_LENGTH}" required></textarea>
          <button class="icomment-button" id="icomment-submit">Submit</button>
        </div>
        <div id="icomment-list" class="icomment-list"></div>
      </div>
    `;

    document.getElementById('icomment-submit').addEventListener('click', () => {
      const nickname = document.getElementById('icomment-nickname').value.trim();
      const email = document.getElementById('icomment-email').value.trim();
      const content = document.getElementById('icomment-content').value.trim();

      if (!validateInput(nickname, email, content)) {
        return;
      }

      fetch(`${API_BASE}/api/comments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          article_url: ARTICLE_URL,
          nickname,
          email,
          content
        })
      })
      .then(res => {
        if (!res.ok) {
          return res.json().then(err => Promise.reject(err));
        }
        return res.json();
      })
      .then(data => {
        document.getElementById('icomment-nickname').value = '';
        document.getElementById('icomment-email').value = '';
        document.getElementById('icomment-content').value = '';
        alert('Comment submitted successfully! It will appear after approval.');
        loadComments();
      })
      .catch(err => {
        const message = err.error || 'Failed to submit comment';
        alert(message + '\n\nPlease try again or refresh the page.');
      });
    });

    loadComments();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
